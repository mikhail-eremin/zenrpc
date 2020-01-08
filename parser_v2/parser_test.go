package parser_v2

import "testing"

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		entryPoint string
		wantErr    bool
	}{
		{
			name:       "Should load packages",
			entryPoint: "/Users/a.simonov/Documents/projects/golang/zenrpc/testdata/catalogue.go",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				entryPoint: tt.entryPoint,
			}
			if err := p.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
