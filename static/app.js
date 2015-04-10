var nad = nad || {};

nad.fmtCmd = function(data) {
  return 'Main.' + [data.Variable, data.Value].join(data.Operator);
};

nad.send = function(ctrl, req) {
  m.request({method: 'POST', url: '/api/v1/nad', data: req})
    .then(function (data) {
      var state = ctrl.model().state;
      if (data.Value === 'On' || data.Value === 'Off') {
        state[data.Variable] = data.Value === 'On';
      }
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
  ctrl.amp = function() {
    nad.send(ctrl, {
      'Variable': 'Model',
      'Operator': '?'
    });
  };
};

nad.console = function(ctrl) {
  var text;
  if (Object.keys(ctrl.model()).length <= 1) {
    text = ['These go to eleven!'];
  } else {
    text = ['sent:     ' + ctrl.model().message,
            'received: ' + ctrl.model().reply];
  }
  return m('pre.console', text.join('\n'));
};

nad.onoff = function(ctrl, options) {
  var isOn = !!ctrl.model().state[options.type];
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

nad.amp = function(ctrl) {
  return m('button[type=button]', {
    class: 'btn btn-default',
    onclick: ctrl.amp
  }, 'Model');
};

nad.error = function(ctrl) {
  var e = ctrl.error();
  var isError = Object.keys(e).length !== 0;
  var text = isError ? e.message + ' (' + e.status + ')' : '';
  var cls = 'alert-danger' + (isError ? '' : ' hidden');
  return m('div.alert',{class: cls, role: 'alert'}, [
    m('strong', 'Error! '), text
  ]);
};

nad.view = function(ctrl) {
  return m('div.container', [
    m('h1', 'NAD Remote'),
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
    m('div.row', {class: 'top-spacing'}, [
      m('div.col-md-4', [nad.source(ctrl)])
    ]),
    m('div.row', {class: 'top-spacing'}, [
      m('div.col-md-2', {class: 'col-md-offset-2'}, [nad.amp(ctrl)])
    ])
  ]);
};

m.module(document.getElementById('nad-remote'), nad);
