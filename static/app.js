var nad = nad || {};

nad.fmtCmd = function(data) {
  return 'Main.' + [data.Variable, data.Value].join(data.Operator);
};

nad.send = function(ctrl, req) {
  m.request({method: 'POST', url: '/api/v1/nad', data: req})
    .then(function (data) {
      ctrl.data = data;
      return data;
    })
    .then(function (data) {
      var msg = nad.fmtCmd(req);
      var reply = nad.fmtCmd(data);
      ctrl.command = {message: msg, reply: reply};
      return data;
    })
    .then(function (data) {
      if (data.Value === 'On' || data.Value === 'Off') {
        ctrl.state[data.Variable] = data.Value === 'On';
      }
    });
};

nad.controller = function() {
  var ctrl = this;
  ctrl.state = {};
  ctrl.power = function() {
    nad.send(ctrl, {
      'Variable': 'Power',
      'Operator': '=',
      'Value': ctrl.state.Power ? 'Off' : 'On'
    });
  };
  ctrl.mute = function() {
    nad.send(ctrl, {
      'Variable': 'Mute',
      'Operator': '=',
      'Value': ctrl.state.Mute ? 'Off' : 'On'
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
  if (ctrl.command) {
      text = ['--> sent:     ' + ctrl.command.message,
              '<-- received: ' + ctrl.command.reply];
  } else {
      text = ['These go to eleven!'];
  }
  return m('pre.console', text.join('\n'));
};

nad.onoff = function(ctrl, options) {
  var isOn = !!ctrl.state[options.type];
  return m('button[type=button]', {
    style: 'width: 100%',
    class: 'btn btn-default btn-lg' + (isOn ? ' active' : ''),
    onclick: options.onclick
  }, [options.icon]);
};

nad.volume = function(ctrl, options) {
  return m('button[type=button]', {
    style: 'width: 100%',
    class: 'btn btn-default btn-lg',
    onclick: options.onclick
  }, [options.icon]);
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
          icon: m('span', {
            class: 'glyphicon glyphicon-off', 'aria-hidden': true
          })
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.onoff(ctrl, {
          onclick: ctrl.mute,
          type: 'Mute',
          icon: m('span', {
            class: 'glyphicon glyphicon-volume-off', 'aria-hidden': true
          })
        })
      ])
    ]),
    m('div.row', [
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.volume(ctrl, {
          onclick: ctrl.volumeUp, type: '+',
          icon: m('span', {
            class: 'glyphicon glyphicon-volume-up', 'aria-hidden': true
          })
        })
      ]),
      m('div.col-md-2', {class: 'top-spacing'}, [
        nad.volume(ctrl, {
          onclick: ctrl.volumeDown, type: '-',
          icon: m('span', {
            class: 'glyphicon glyphicon-volume-down', 'aria-hidden': true
          })
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
