var nad = nad || {};

nad.bindKeys = function(ctrl) {
  Mousetrap.bind('p', ctrl.power);
  Mousetrap.bind('m', ctrl.mute);
  Mousetrap.bind('s', ctrl.speakerA);
  Mousetrap.bind('i', ctrl.amp);
  Mousetrap.bind('+', ctrl.volumeUp);
  Mousetrap.bind('-', ctrl.volumeDown);
};

nad.fmtCmd = function(data) {
  return 'Main.' + [data.Variable, data.Value].join(data.Operator);
};

nad.send = function(ctrl, req) {
  m.request({method: 'POST', url: '/api/v1/nad', data: req})
    .then(function (data) {
      var state = ctrl.model().state;
      state[data.Variable] = data.Value === 'On' || data.Value === 'Off' ?
        data.Value === 'On' : data.Value;
      ctrl.error({});
      ctrl.model({
        message: nad.fmtCmd(req),
        reply: nad.fmtCmd(data),
        state: state
      });
    }, ctrl.error);
};

nad.controller = function() {
  var ctrl = this;
  ctrl.error = m.prop({});
  ctrl.model = m.prop({state: {}});
  ctrl.power = function() {
    nad.send(ctrl, {
      'Variable': 'Power',
      'Operator': '=',
      'Value': ctrl.model().state.Power ? 'Off' : 'On'
    });
  };
  ctrl.mute = function() {
    nad.send(ctrl, {
      'Variable': 'Mute',
      'Operator': '=',
      'Value': ctrl.model().state.Mute ? 'Off' : 'On'
    });
  };
  ctrl.volumeUp = function() {
    nad.send(ctrl, {
      'Variable': 'Volume',
      'Operator': '+',
    });
  };
  ctrl.volumeDown = function() {
    nad.send(ctrl, {
      'Variable': 'Volume',
      'Operator': '-',
    });
  };
  ctrl.source = function(value) {
    nad.send(ctrl, {
      'Variable': 'Source',
      'Operator': '=',
      'Value': value
    });
  };
  ctrl.refreshSource = function() {
    nad.send(ctrl, {
      'Variable': 'Source',
      'Operator': '?'
    });
  };
  ctrl.amp = function() {
    nad.send(ctrl, {
      'Variable': 'Model',
      'Operator': '?'
    });
  };
  ctrl.speakerA = function() {
    // Assume that initial state is on
    var spkrA = ctrl.model().state.SpeakerA;
    var isOn = _.isUndefined(spkrA) ? true : spkrA;
    nad.send(ctrl, {
      'Variable': 'SpeakerA',
      'Operator': '=',
      'Value': isOn ? 'Off' : 'On'
    });
  };
  nad.bindKeys(ctrl);
};

nad.console = function(ctrl) {
  var text;
  if (_.isEmpty(ctrl.model().state)) {
    text = ['These go to eleven!'];
  } else {
    text = ['sent:     ' + ctrl.model().message,
            'received: ' + ctrl.model().reply];
  }
  return m('pre.console', text.join('\n'));
};

nad.onoff = function(ctrl, options) {
  var state = ctrl.model().state[options.type];
  var isOn = !!state;
  if (_.isUndefined(state) && !_.isUndefined(options.initialState)) {
    isOn = options.initialState;
  }
  var active = options.invert ? !isOn : isOn;
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg' + (active ? ' active' : ''),
    onclick: options.onclick
  }, options.icon);
};

nad.volume = function(ctrl, options) {
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg',
    onclick: options.onclick
  }, options.icon);
};

nad.source = function(ctrl) {
  var sources = ['CD', 'Tuner', 'Video', 'Disc/MDC', 'Tape2', 'Aux'];
  var model = ctrl.model();
  return m('select.form-control', {
    onchange: m.withAttr('value', ctrl.source)
  }, _.map(sources, function(src) {
    var val = src.toUpperCase();
    var selected = model.state.Source === val ? 'selected' : '';
    return m('option', {'value': val, 'selected': selected}, src);
  }));
};

nad.refreshSource = function(ctrl, options) {
  return m('button[type=button]', {
    class: 'btn btn-default',
    onclick: ctrl.refreshSource
  }, options.icon);
};

nad.amp = function(ctrl, options) {
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg',
    onclick: ctrl.amp
  }, options.icon);
};

nad.error = function(ctrl) {
  var e = ctrl.error();
  var isError = !_.isEmpty(e);
  var text = isError ? e.message + ' (' + e.status + ')' : '';
  var cls = 'alert-danger' + (isError ? '' : ' hidden');
  return m('div.alert', {class: cls, role: 'alert'}, [
    m('strong', 'Error! '), text
  ]);
};

nad.view = function(ctrl) {
  return m('div.container', [
    m('div.row', [
      m('div.col-md-4', m('h1', [
        m('span', {class: 'glyphicon glyphicon-signal'}), ' amp remote'
      ]))
    ]),
    m('div.row', [
      m('div.col-md-4', nad.error(ctrl))
    ]),
    m('div.row', [
      m('div.col-md-4', [
        nad.console(ctrl)
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff(ctrl, {
          onclick: ctrl.power,
          type: 'Power',
          icon: m('span', {class: 'glyphicon glyphicon-off'})
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff(ctrl, {
          onclick: ctrl.mute,
          type: 'Mute',
          icon: m('span', {class: 'glyphicon glyphicon-volume-off'})
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.volume(ctrl, {
          onclick: ctrl.volumeUp,
          icon: m('span', {class: 'glyphicon glyphicon-volume-up'})
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.volume(ctrl, {
          onclick: ctrl.volumeDown,
          icon: m('span', {class: 'glyphicon glyphicon-volume-down'})
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff(ctrl, {
          onclick: ctrl.speakerA,
          type: 'SpeakerA',
          icon: m('span', {class: 'glyphicon glyphicon-headphones'}),
          initialState: true,
          invert: true
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.amp(ctrl, {
          icon: m('span', {class: 'glyphicon glyphicon-info-sign'})
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, nad.source(ctrl)),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.refreshSource(ctrl, {
          icon: m('span', {class: 'glyphicon glyphicon-refresh'})
        })
      ])
    ])
  ]);
};

m.mount(document.getElementById('nad-remote'), nad);
