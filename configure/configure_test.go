package configure

import (
	"bytes"
	"testing"
)

func TestFromStrings(t *testing.T) {
	type args struct {
		discordWH string
		slackWH   string
		msgTmpl   string
	}
	tests := []struct {
		name    string
		args    args
		testMsg string
		wantMsg string
		wantErr bool
	}{
		{
			name: "valid options",
			args: args{
				discordWH: "http://d.example.com",
				slackWH:   "http://s.example.com",
				msgTmpl:   "TEST {{.Message}} TEST",
			},
			testMsg: "hej svejs",
			wantMsg: "TEST hej svejs TEST",
		},
		{
			name: "optional options",
			args: args{},
		},
		{
			name: "bad template",
			args: args{
				msgTmpl: "TEST {{.WrongThing}} TEST",
			},
			wantErr: true,
		},
		{
			name: "invalid webhook URL 1",
			args: args{
				discordWH: "file://example.com",
				msgTmpl:   "TEST {{.Message}} TEST",
			},
			wantErr: true,
		},
		{
			name: "invalid webhook URL 2",
			args: args{
				slackWH: "ftp-example-com",
				msgTmpl: "TEST {{.Message}} TEST",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := FromStrings(tt.args.discordWH, tt.args.slackWH, tt.args.msgTmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.args.discordWH != "" && gotCfg.DiscordHookURL.String() != tt.args.discordWH {
				t.Errorf("FromStrings() Discord hook URL = %s, want %s", gotCfg.DiscordHookURL, tt.args.discordWH)
			}
			if tt.args.slackWH != "" && gotCfg.SlackHookURL.String() != tt.args.slackWH {
				t.Errorf("FromStrings() Slack hook URL = %s, want %s", gotCfg.SlackHookURL, tt.args.slackWH)
			}

			if tt.wantMsg != "" {
				var buf bytes.Buffer
				err = gotCfg.MessageTemplate.Execute(&buf, struct{ Message string }{tt.testMsg})
				if err != nil {
					t.Errorf("FromStrings() MessageTemplate can't render %s: %v", tt.testMsg, err)
				}
				gotMsg := buf.String()
				if gotMsg != tt.wantMsg {
					t.Errorf("FromStrings() MessageTemplate rendered %s, want %s", gotMsg, tt.wantMsg)
				}
			}
		})
	}
}
