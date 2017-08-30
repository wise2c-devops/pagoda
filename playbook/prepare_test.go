package playbook

import "testing"

func TestGetTemplatePath(t *testing.T) {
	tp, err := getTemplatePath("good-playbook")
	if err != nil {
		t.Error(err)
	}
	t.Log(tp)

	_, err = getTemplatePath("bad-playbook")
	if err == nil {
		t.Error("bad-playbook should have error")
	}
}
