package appidentity

const (
	// AppUserModelID is shared by every current-generation Windows process,
	// shortcut, and toast that users perceive as Temper. Keep it
	// version-independent across upgrades. It matches the MSIX package
	// identity (Temper.Cowork.Desktop) so taskbar grouping and toasts align
	// with the packaged app.
	AppUserModelID = "Temper.Cowork.Desktop"

	// legacyTauriAppUserModelID was written to Windows shortcuts by Reasonix
	// Desktop 0.53. Keep the current identity distinct so separately installed
	// Reasonix and Temper generations do not merge into one taskbar group.
	legacyTauriAppUserModelID = "dev.reasonix.desktop"
)
