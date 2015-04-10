var nad = nad || {};

nad.fmtCmd = function(data) {
  return 'Main.' + [data.Variable, data.Value].join(data.Operator);
};

nad.send = function(ctrl, req) {
  m.request({method: 'POST', url: '/api/v1/nad', data: req})
    .then(function (data) {
      ctrl.model({message: nad.fmtCmd(req), reply: nad.fmtCmd(data)});
      return data;
    })
    .then(function (data) {
      if (data.Value === 'On' || data.Value === 'Off') {
        var state = ctrl.state();
        state[data.Variable] = data.Value === 'On';
        ctrl.state(state);
      }
    });
};

nad.controller = function() {
  var ctrl = this;
  ctrl.state = m.prop({});
  ctrl.model = m.prop({});
  ctrl.power = function() {
    nad.send(ctrl, {
      'Variable': 'Power',
      'Operator': '=',
      'Value': ctrl.state().Power ? 'Off' : 'On'
    });
  };
  ctrl.mute = function() {
    nad.send(ctrl, {
      'Variable': 'Mute',
      'Operator': '=',
      'Value': ctrl.state().Mute ? 'Off' : 'On'
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
};

nad.console = function(ctrl) {
  var text;
  if (Object.keys(ctrl.model()).length === 0) {
    text = ['These go to eleven!'];
  } else {
    text = ['sent:     ' + ctrl.model().message,
            'received: ' + ctrl.model().reply];
  }
  return m('pre.console', text.join('\n'));
};

nad.onoff = function(ctrl, options) {
  var isOn = !!ctrl.state()[options.type];
  return m('button[type=button]', {
    class: 'btn btn-default btn-lg' + (isOn ? ' active' : ''),
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
  return m('select.form-control', {
    onchange: m.withAttr('value', ctrl.source)
  }, [
    m('option[value=CD]', 'CD'),
    m('option[value=TUNER]', 'Tuner'),
    m('option[value=VIDEO]', 'Video'),
    m('option[value=DISC/MDC]', 'Disc/MDC'),
    m('option[value=TAPE2]', 'Tape2'),
    m('option[value=AUX]', 'Aux')
  ]);
};

nad.view = function(ctrl) {
  return m('div.container', [
    m('h1', 'NAD Remote'),
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
    m('div.row', {class: 'top-spacing'}, [
      m('div.col-md-4', [
        nad.source(ctrl, {onclick: ctrl.source})
      ])
    ])
  ]);
};

m.module(document.getElementById('nad-remote'), nad);
