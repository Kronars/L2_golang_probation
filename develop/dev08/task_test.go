package main_test

import "testing"
import "dev08"

func Test_execute(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
		want1 bool
	}{
		{"pwd", []string{"pwd"}, `c:\Users\Senya\Prog_2\Projects\WB probation\L2\tasks\develop\dev08`, true},
		{"cd", []string{"cd", ".."}, "", true},
		{"echo", []string{"echo", "<sarcasm> это хорошие юнит тесты </sarcasm>"}, "<sarcasm> это хорошие юнит тесты </sarcasm>", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := main.Execute(tt.parts)
			if got != tt.want {
				t.Errorf("execute() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("execute() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
