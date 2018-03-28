package nad

import "testing"

func TestSendCmd(t *testing.T) {
	nad := NewTestClient()
	_, err := nad.SendCmd(Cmd{Variable: "foo"})
	if err == nil {
		t.Error("Expected error")
	}
	actual, err := nad.SendCmd(Cmd{Variable: "Model", Operator: "?"})
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Model=C356BEE"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestModel(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Model()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Model=C356BEE"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestEnable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.enable("Power", true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Power=On"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
	actual, err = nad.enable("Power", false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Power=Off"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestMuteEnable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Mute(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Mute=On"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestMuteDisable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Mute(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Mute=Off"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestPowerEnable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Power(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Power=On"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestPowerDisable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Power(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Power=Off"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSource(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Source("CD")
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Source=CD"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerAEnable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.SpeakerA(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerA=On"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerADisable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.SpeakerA(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerA=Off"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerBEnable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.SpeakerB(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerB=On"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerBDisable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.SpeakerB(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerB=Off"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestTape1Enable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Tape1(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Tape1=On"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestTape1Disable(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.Tape1(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Tape1=Off"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestVolumeUp(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.VolumeUp()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Volume+"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestVolumeDown(t *testing.T) {
	nad := NewTestClient()
	actual, err := nad.VolumeDown()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Volume-"; actual.String() != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestEnableVolume(t *testing.T) {
	nad := NewTestClient()
	nad.EnableVolume = false
	if _, err := nad.VolumeUp(); err == nil {
		t.Error("Expected error")
	}
	if _, err := nad.VolumeDown(); err == nil {
		t.Error("Expected error")
	}
	volumeUp := Cmd{Variable: "Volume", Operator: "+"}
	if _, err := nad.SendCmd(volumeUp); err == nil {
		t.Error("Expected error")
	}
	volumeDown := Cmd{Variable: "Volume", Operator: "-"}
	if _, err := nad.SendCmd(volumeDown); err == nil {
		t.Error("Expected error")
	}
}

func TestEnsureOpen(t *testing.T) {
	nad := NewTestClient()
	if _, err := nad.Power(true); err != nil {
		t.Fatal(err)
	}
	if want := "/dev/realfoo"; nad.device.realname != want {
		t.Errorf("want %s, got %s", want, nad.device.realname)
	}
	// Device symlink changes
	nad.device.evalSymlinks = func(string) (string, error) { return "/dev/realbar", nil }
	// Sending command updates device realname
	if _, err := nad.Power(true); err != nil {
		t.Fatal(err)
	}
	if want := "/dev/realbar"; nad.device.realname != want {
		t.Errorf("want %s, got %s", want, nad.device.realname)
	}
}
