package main

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
)

func TestCgoCallbackTemplatesUsePointerParameters(t *testing.T) {
	t.Parallel()

	params := rpcParams{
		ServiceName: "State",
		MethodName:  "GetState",
		RequestType: "GetStateRequest",
		CgoBindings: true,
	}

	tests := []struct {
		name    string
		tmpl    *template.Template
		want    []string
		notWant []string
	}{
		{
			name: "unary",
			tmpl: syncTemplate,
			want: []string{
				"callback *C.CCallback",
				"WrapCallbackCgo(*callback)",
			},
			notWant: []string{
				"callback C.CCallback",
				"if callback == nil",
			},
		},
		{
			name: "server stream",
			tmpl: readStreamTemplate,
			want: []string{
				"rStream *C.CRecvStream",
				"WrapRecvStreamCgo(*rStream)",
			},
			notWant: []string{
				"rStream C.CRecvStream",
				"if rStream == nil",
			},
		},
		{
			name: "bidirectional stream",
			tmpl: biStreamTemplate,
			want: []string{
				"rStream *C.CRecvStream",
				"return C.uintptr_t(0)",
				"WrapRecvStreamCgo(*rStream)",
			},
			notWant: []string{
				"rStream C.CRecvStream",
				"if rStream == nil",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			if err := test.tmpl.Execute(&output, params); err != nil {
				t.Fatalf("execute template: %v", err)
			}

			generated := output.String()
			for _, want := range test.want {
				if !strings.Contains(generated, want) {
					t.Errorf("generated output does not contain %q", want)
				}
			}
			for _, notWant := range test.notWant {
				if strings.Contains(generated, notWant) {
					t.Errorf("generated output unexpectedly contains %q", notWant)
				}
			}
		})
	}
}
