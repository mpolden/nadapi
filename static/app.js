var nad = nad || {};

nad.state = {
  // Initial amplifier state
  amp: {Power: false, Mute: false, Speakera: true},
  message: {},
  error: {},
  helpVisible: false,
  power: function() {
    nad.send({
      Variable: 'Power',
      Operator: '=',
      Value: !nad.state.amp.Power
    });
  },
  mute: function() {
    nad.send({
      Variable: 'Mute',
      Operator: '=',
      Value: !nad.state.amp.Mute
    });
  },
  volumeUp: function() {
    nad.send({
      Variable: 'Volume',
      Operator: '+',
    });
  },
  volumeDown: function() {
    nad.send({
      Variable: 'Volume',
      Operator: '-',
    });
  },
  source: function(value) {
    nad.send({
      Variable: 'Source',
      Operator: '=',
      Value: value
    });
  },
  ampModel: function() {
    nad.send({
      Variable: 'Model',
      Operator: '?'
    });
  },
  speakerA: function() {
    nad.send({
      Variable: 'Speakera',
      Operator: '=',
      Value: !nad.state.amp.Speakera
    });
  },
  reload: function(refresh) {
    var data = refresh ? {refresh: true} : {};
    nad.get('Power', data);
    nad.get('Mute', data);
    nad.get('Source', data);
    nad.get('Speakera', data);
  },
  toggleHelp: function() {
    nad.state.helpVisible = !nad.state.helpVisible;
    m.redraw();
  }
};

nad.bindKeys = function() {
  nad.keyBindings.forEach(function (kb) {
    Mousetrap.bind(kb.key, nad.state[kb.callback]);
  });
};

nad.keyBindings = [
  {key: 'p', callback: 'power', description: 'Toggle power'},
  {key: 'm', callback: 'mute', description: 'Toggle mute'},
  {key: 's', callback: 'speakerA', description: 'Toggle headphones'},
  {key: 'i', callback: 'amp', description: 'Get amplifier model'},
  {key: 'r', callback: 'reload', description: 'Reload state from amplifier'},
  {key: '+', callback: 'volumeUp', description: 'Increase volume'},
  {key: '-', callback: 'volumeDown', description: 'Decrease volume'},
  {key: 'h', callback: 'toggleHelp',
   description: 'Togge list of keyboard shortcuts'}
];

nad.fmtCmd = function(data) {
  return 'Main.' + [data.Variable, data.Value].join(data.Operator);
};

nad.fromValue = function(v) {
  return v === 'On' || v === 'Off' ? v === 'On' : v;
};

nad.toValue = function(v) {
  if (v === true || v === false) {
    return v ? 'On' : 'Off';
  }
  return v;
};

nad.get = function(variable, data) {
  var url = '/api/v1/nad/state/' + variable;
  m.request({method: 'GET', url: url, data: data})
    .then(function (data) {
      var amp = nad.state.amp;
      amp[data.Variable] = nad.fromValue(data.Value);
    }, function (data) {
      nad.state.error = data;
    });
};

nad.send = function(req) {
  req.Value = nad.toValue(req.Value);
  m.request({method: 'POST', url: '/api/v1/nad', data: req})
    .then(function (data) {
      var amp = nad.state.amp;
      amp[data.Variable] = nad.fromValue(data.Value);
      nad.state.error = {};
      nad.state.message = {request: req, reply: data};
    }, function (data) {
      nad.state.error = data;
    });
};

nad.console = function() {
  var text;
  if (Object.keys(nad.state.message).length === 0) {
    text = ['These go to eleven!'];
  } else {
    text = ['sent:     ' + nad.fmtCmd(nad.state.message.request),
            'received: ' + nad.fmtCmd(nad.state.message.reply)];
  }
  return m('pre.console', text.join('\n'));
};

nad.onoff = function(options) {
  var amp = nad.state.amp;
  if (!amp.hasOwnProperty(options.type)) {
    throw 'Unknown type: ' + options.type;
  }
  var isOn = amp[options.type];
  var active = options.invert ? !isOn : isOn;
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg' + (active ? ' active' : ''),
    onclick: options.onclick
  }, options.icon);
};

nad.volume = function(options) {
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg',
    onclick: options.onclick
  }, options.icon);
};

nad.source = function() {
  var sources = ['CD', 'Tuner', 'Video', 'Disc/MDC', 'Tape2', 'Aux'];
  var amp = nad.state.amp;
  return m('select.form-control', {
    onchange: m.withAttr('value', nad.state.source)
  }, sources.map(function(src) {
    var val = src.toUpperCase();
    var selected = amp.Source === val ? 'selected' : '';
    return m('option', {value: val, selected: selected}, src);
  }));
};

nad.reloadState = function(options) {
  return m('button[type=button]', {
    class: 'btn btn-default',
    onclick: nad.state.reload
  }, options.icon);
};

nad.amp = function(options) {
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg',
    onclick: nad.state.ampModel
  }, options.icon);
};

nad.error = function() {
  var e = nad.state.error;
  var isError = Object.keys(e).length !== 0;
  var text = isError ? e.message + ' (' + e.status + ')' : '';
  var cls = 'alert-danger' + (isError ? '' : ' hidden');
  return m('div.alert', {class: cls, role: 'alert'}, [
    m('strong', 'Error! '), text
  ]);
};

nad.help = function() {
  if (!nad.state.helpVisible) {
    return m('p.text-muted', 'Tip: Press ', m('code', 'h'),
             ' to display keyboard shortcuts');
  }
  var rows = nad.keyBindings.map(function (kb) {
    return m('tr', [
      m('td', m('center', m('code', kb.key))),
      m('td', kb.description)
    ]);
  });
  return m('table.table',
           m('thead', m('tr', [
             m('th', 'Key binding'),
             m('th', 'Description')
           ])),
           m('tbody', rows)
          );
};

nad.oninit = nad.bindKeys;

nad.view = function() {
  return m('div.container', [
    m('div.row', [
      m('div.col-md-4', m('h1', [
        m('span', {class: 'glyphicon glyphicon-signal'}), ' amp remote'
      ]))
    ]),
    m('div.row', [
      m('div.col-md-4', nad.error())
    ]),
    m('div.row', [
      m('div.col-md-4', [
        nad.console()
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff({
          onclick: nad.state.power,
          type: 'Power',
          icon: m('span', {class: 'glyphicon glyphicon-off'})
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff({
          onclick: nad.state.mute,
          type: 'Mute',
          icon: m('span', {class: 'glyphicon glyphicon-volume-off'})
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.volume({
          onclick: nad.state.volumeUp,
          icon: m('span', {class: 'glyphicon glyphicon-volume-up'})
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.volume({
          onclick: nad.state.volumeDown,
          icon: m('span', {class: 'glyphicon glyphicon-volume-down'})
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff({
          onclick: nad.state.speakerA,
          type: 'Speakera',
          icon: m('span', {class: 'glyphicon glyphicon-headphones'}),
          invert: true
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.amp({
          icon: m('span', {class: 'glyphicon glyphicon-info-sign'})
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, nad.source()),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.reloadState({
          icon: m('span', {class: 'glyphicon glyphicon-refresh'})
        })
      ])
    ]),
    m('div.row', m('div.col-md-4', {class: 'top-spacing'}, nad.help()))
  ]);
};

m.mount(document.getElementById('nad-remote'), nad);
