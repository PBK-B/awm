package workspace

import "awmcli/internal/gitutil"

func RefreshLockEntry(name string, dep Dependency) error {
	l, err := ReadLock()
	if err != nil {
		if Missing(err) {
			root := "workspace"
			if m, merr := ReadMetadata(); merr == nil {
				root = m.Name
			}
			l = DefaultLock(root)
		} else {
			return err
		}
	}
	l.Dependencies[name] = LockEntry{URL: dep.URL, Path: dep.Path, Manager: dep.Manager, Role: dep.Role, RequestedBranch: dep.Branch, ResolvedBranch: gitutil.Branch(dep.Path), Commit: gitutil.Commit(dep.Path), Remote: "origin", Dirty: gitutil.Dirty(dep.Path), Installed: Exists(dep.Path), LastCheckedAt: Now()}
	return WriteLock(l)
}
