package powerprofiles

import (
	"bytes"

	"git.sr.ht/~rockorager/vaxis"

	"github.com/nekorg/pawbar/internal/config"
	"github.com/nekorg/pawbar/internal/modules"
)

type powerProfileModule struct {
	receive chan bool
	send    chan modules.Event

	opts        Options
	initialOpts Options
}

func (mod *powerProfileModule) Dependencies() []string {
	return nil
}

func (mod *powerProfileModule) Channels() (<-chan bool, chan<- modules.Event) {
	return mod.receive, mod.send
}

func (mod *powerProfileModule) Name() string {
	return "powerprofiles"
}

func (mod *powerProfileModule) Run() (<-chan bool, chan<- modules.Event, error) {
	mod.receive = make(chan bool)
	mod.send = make(chan modules.Event)
	mod.initialOpts = mod.opts

	go func() {
		for {
			select {
			case e := <-mod.send:
				switch ev := e.VaxisEvent.(type) {
				case vaxis.Mouse:
					if ev.EventType != vaxis.EventPress {
						break
					}
					btn := config.ButtonName(ev)

					if mod.opts.OnClick.Dispatch(btn, &mod.initialOpts, &mod.opts) {
						mod.receive <- true
					}

					switch ev.Button {
					case vaxis.MouseRightButton:
						// go power.LaunchMenu(ev.XPixel/2, ev.YPixel/2)
					}

				case modules.FocusIn:
					if mod.opts.OnClick.HoverIn(&mod.opts) {
						mod.receive <- true
					}

				case modules.FocusOut:
					if mod.opts.OnClick.HoverOut(&mod.opts) {
						mod.receive <- true
					}
				}
			}
		}
	}()

	return mod.receive, mod.send, nil
}

func (mod *powerProfileModule) Render() []modules.EventCell {
	var s vaxis.Style
	s.Foreground = mod.opts.Fg.Go()
	s.Background = mod.opts.Bg.Go()

	var buf bytes.Buffer
	_ = mod.opts.Format.Execute(&buf, nil)

	chars := vaxis.Characters(buf.String())
	r := make([]modules.EventCell, len(chars))

	for i, ch := range chars {
		r[i] = modules.EventCell{
			C: vaxis.Cell{
				Character: ch,
				Style:     s,
			},
			Metadata:   "",
			Mod:        mod,
			MouseShape: mod.opts.Cursor.Go(),
		}
	}
	return r
}
