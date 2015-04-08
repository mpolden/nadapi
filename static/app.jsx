var nad = nad || {};

nad.sendToAmp = function(req, success) {
  $.ajax({
    url: '/api/v1/nad',
    type: 'POST',
    data: JSON.stringify(req),
    dataType: 'json',
    processData: false,
    success: success,
    error: function(xhr, status, err) {
      console.error(xhr, status, err.toString());
    }.bind(this)
  });
};

nad.Console = React.createClass({
  render: function() {
    var text;
    if (this.props.command) {
      text = ['--> sent:     ' + this.props.command.message,
              '<-- received: ' + this.props.command.reply];
    } else {
      text = ['These go to eleven!'];
    }
    return (
      <pre>
        {text.join('\n')}
      </pre>
    );
  }
});

nad.Source = React.createClass({
  handleChange: function(event) {
    var req = {
      'Variable': 'Source',
      'Operator': '=',
      'Value': event.target.value
    };
    nad.sendToAmp(req, function(resp) {
        this.props.onUpdate(req, resp);
    }.bind(this));
  },
  render: function() {
    return (
      <select className='form-control' onChange={this.handleChange}>
        <option value='CD'>CD</option>
        <option value='TUNER'>Tuner</option>
        <option value='VIDEO'>Video</option>
        <option value='DISC/MDC'>Disc/MDC</option>
        <option value='TAPE2'>Tape2</option>
        <option value='AUX'>Aux</option>
      </select>
    );
  }
});

nad.Volume = React.createClass({
  handleClick: function() {
    var req = {
      'Variable': 'Volume',
      'Operator': this.props.type
    };
    nad.sendToAmp(req, function(resp) {
        this.props.onUpdate(req, resp);
    }.bind(this));
  },
  shouldComponentUpdate: function(nextProps) {
    return nextProps.type === nextProps.data.Operator;
  },
  render: function() {
    return (
      <button type='button' className='btn btn-default' style={{width: '100%'}}
              onClick={this.handleClick}>
        Volume {this.props.type}
      </button>
    );
  }
});

nad.OnOffButton = React.createClass({
  handleClick: function() {
    var req = {
      'Variable': this.props.type,
      'Operator': '=',
      'Value': this.props.data.Value === 'On' ? 'Off' : 'On'
    };
    nad.sendToAmp(req, function(resp) {
        this.props.onUpdate(req, resp);
    }.bind(this));
  },
  shouldComponentUpdate: function(nextProps) {
    return nextProps.type === nextProps.data.Variable;
  },
  render: function() {
    var isOn = this.props.data.Value === 'On';
    var text = isOn ? 'off' : 'on';
    return (
      <button type='button' style={{width: '100%'}}
        className={isOn ? 'btn btn-success' : 'btn btn-default'}
        onClick={this.handleClick}>
        {this.props.type} {text}
      </button>
    );
  }
});

nad.Remote = React.createClass({
  fmtCmd: function(data) {
    return 'Main.' + [data.Variable, data.Value].join(data.Operator);
  },
  getInitialState: function() {
    return {data: {}};
  },
  onUpdate: function(req, resp) {
    var msg = this.fmtCmd(req);
    var reply = this.fmtCmd(resp);
    this.setState({data: resp, command: {message: msg, reply: reply}});
  },
  render: function() {
    return (
      <div className='container'>
        <h1>NAD Remote</h1>
        <div className='row'>
          <div className='col-md-4'>
            <nad.Console command={this.state.command} />
          </div>
        </div>
        <div className='row top-spacing'>
        <div className='col-md-2'><nad.OnOffButton type='Power'
             onUpdate={this.onUpdate} data={this.state.data} />
          </div>
        <div className='col-md-2'><nad.OnOffButton type='Mute'
             onUpdate={this.onUpdate} data={this.state.data} />
        </div>
        </div>
        <div className='row top-spacing'>
        <div className='col-md-2'><nad.Volume type='+'
             onUpdate={this.onUpdate} data={this.state.data} />
          </div>
        <div className='col-md-2'><nad.Volume type='-'
             onUpdate={this.onUpdate} data={this.state.data} />
        </div>
        </div>
        <div className='row top-spacing'>
          <div className='col-md-4'>
            <nad.Source onUpdate={this.onUpdate} data={this.state.data} />
          </div>
        </div>
      </div>
    );
  }
});

React.render(
  <nad.Remote />,
  document.getElementById('nad-remote')
);
