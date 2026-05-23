# Schema reference

**Source:** [SafeInCloud](https://www.safe-in-cloud.com) macOS client, version
**v25.3.5**, sourced on **2026-05-23** from
`/Applications/Safe.app/Contents/Resources/`.

| File | Purpose |
|---|---|
| `database.xml` | Localized template definition (`@string/...` placeholders, no user data). |
| `database_demo.xml` | Same templates plus sample card data used by the in-app demo. |

These are ground truth for the Go XML structs in `unmarshal.go`. When SafeInCloud
ships a new format version, refresh these files and reconcile any new attributes
or elements with the Go structs.
