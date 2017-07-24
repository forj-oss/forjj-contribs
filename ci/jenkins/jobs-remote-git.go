package main

func (t *GitStruct)SetFrom(d *GitStruct) bool {
	if t == nil {
		return false
	}
	SetIfSet(&t.RemoteUrl, d.RemoteUrl)
	return true
}

func (t *GitStruct)UpdateFrom(d *GitStruct) bool {
	if t == nil {
		return false
	}
	SetIfSet(&t.RemoteUrl, d.RemoteUrl)
	return true
}

func (t *GitStruct)GetUpstream() string {
	return t.RemoteUrl
}
