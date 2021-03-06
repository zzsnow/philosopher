package rep

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"philosopher/lib/sys"
)

// Serialize converts the whole structure to a gob file
func (evi *Evidence) Serialize() {

	b, e := msgpack.Marshal(&evi)
	if e != nil {
		logrus.Fatal("cannot marshal file:", e)
	}

	e = ioutil.WriteFile(sys.EvBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize data:", e)
	}

	return
}

// SerializeGranular converts the whole structure into sevral small gob files
func (evi *Evidence) SerializeGranular() {

	// create EV Parameters
	SerializeEVParameters(evi)

	// create EV PSM
	SerializeEVPSM(evi)

	// create EV Ion
	SerializeEVIon(evi)

	// create EV Peptides
	SerializeEVPeptides(evi)

	// create EV Ion
	SerializeEVProteins(evi)

	// create EV Mods
	SerializeEVMods(evi)

	// create EV Modifications
	SerializeEVModifications(evi)

	// create EV Combined
	SerializeEVCombined(evi)

	return
}

// SerializeEVParameters creates an ev serial with Parameter data
func SerializeEVParameters(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.Parameters)
	if e != nil {
		logrus.Trace("Cannot marshal Parameters data:", e)
	}

	e = ioutil.WriteFile(sys.EvParameterBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Parameters data:", e)
	}

	return
}

// SerializeEVPSM creates an ev serial with Evidence data
func SerializeEVPSM(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.PSM)
	if e != nil {
		logrus.Trace("Cannot marshal PSM data:", e)
	}

	e = ioutil.WriteFile(sys.EvPSMBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize PSM data:", e)
	}

	return
}

// SerializeEVIon creates an ev serial with Evidence data
func SerializeEVIon(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.Ions)
	if e != nil {
		logrus.Trace("Cannot marshal Ions data:", e)
	}

	e = ioutil.WriteFile(sys.EvIonBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Ions data:", e)
	}

	return
}

// SerializeEVPeptides creates an ev serial with Evidence data
func SerializeEVPeptides(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.Peptides)
	if e != nil {
		logrus.Trace("Cannot marshal Peptides data:", e)
	}

	e = ioutil.WriteFile(sys.EvPeptideBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Peptides data:", e)
	}

	return
}

// SerializeEVProteins creates an ev serial with Evidence data
func SerializeEVProteins(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.Proteins)
	if e != nil {
		logrus.Trace("Cannot marshal Proteins data:", e)
	}

	e = ioutil.WriteFile(sys.EvProteinBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Proteins data:", e)
	}

	return
}

// SerializeEVMods creates an ev serial with Evidence data
func SerializeEVMods(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.Mods)
	if e != nil {
		logrus.Trace("Cannot marshal Modifications data:", e)
	}

	e = ioutil.WriteFile(sys.EvModificationsBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Modifications data:", e)
	}

	return
}

// SerializeEVModifications creates an ev serial with Evidence data
func SerializeEVModifications(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.Modifications)
	if e != nil {
		logrus.Trace("Cannot marshal data:", e)
	}

	e = ioutil.WriteFile(sys.EvModificationsEvBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize data:", e)
	}

	return
}

// SerializeEVCombined creates an ev serial with Evidence data
func SerializeEVCombined(evi *Evidence) {

	b, e := msgpack.Marshal(&evi.CombinedProtein)
	if e != nil {
		logrus.Trace("Cannot marshal data:", e)
	}

	e = ioutil.WriteFile(sys.EvCombinedBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize data:", e)
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (evi *Evidence) Restore() {

	b, e := ioutil.ReadFile(sys.EvBin())
	if e != nil {
		logrus.Trace("Cannot marshal data:", e)
	}

	e = msgpack.Unmarshal(b, &e)
	if e != nil {
		logrus.Trace("Cannot serialize data:", e)
	}

	return
}

// RestoreGranular reads philosopher results files and restore the data sctructure
func (evi *Evidence) RestoreGranular() {

	// Parameters
	RestoreEVParameters(evi)

	// PSM
	RestoreEVPSM(evi)

	// Ion
	RestoreEVIon(evi)
	// Peptide

	RestoreEVPeptide(evi)

	// Protein
	RestoreEVProtein(evi)

	// Mods
	RestoreEVMods(evi)

	// Modifications
	RestoreEVModifications(evi)

	// Combined
	RestoreEVCombined(evi)

	return
}

// RestoreEVParameters restores Ev PSM data
func RestoreEVParameters(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvParameterBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Parameters)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}
	return
}

// RestoreEVPSM restores Ev PSM data
func RestoreEVPSM(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvPSMBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.PSM)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}
	return
}

// RestoreEVIon restores Ev Ion data
func RestoreEVIon(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvIonBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Ions)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVPeptide restores Ev Ion data
func RestoreEVPeptide(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvPeptideBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Peptides)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVProtein restores Ev Protein data
func RestoreEVProtein(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvProteinBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Proteins)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVMods restores Ev Mods data
func RestoreEVMods(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvModificationsBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Mods)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVModifications restores Ev Mods data
func RestoreEVModifications(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvModificationsEvBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Modifications)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVCombined restores Ev Mods data
func RestoreEVCombined(evi *Evidence) {

	b, e := ioutil.ReadFile(sys.EvCombinedBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.CombinedProtein)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreGranularWithPath reads philosopher results files and restore the data sctructure
func (evi *Evidence) RestoreGranularWithPath(p string) {

	// Parameters
	RestoreEVParametersWithPath(evi, p)

	// PSM
	RestoreEVPSMWithPath(evi, p)

	// Ion
	RestoreEVIonWithPath(evi, p)

	// Peptide
	RestoreEVPeptideWithPath(evi, p)

	// Protein
	RestoreEVProteinWithPath(evi, p)

	// Mods
	RestoreEVModsWithPath(evi, p)

	// Modifications
	RestoreEVModificationsWithPath(evi, p)

	// Combined
	RestoreEVCombinedWithPath(evi, p)

	return
}

// RestoreEVParametersWithPath restores Ev PSM data
func RestoreEVParametersWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvParameterBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Parameters)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVPSMWithPath restores Ev PSM data
func RestoreEVPSMWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPSMBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.PSM)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVIonWithPath restores Ev Ion data
func RestoreEVIonWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvIonBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Ions)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVPeptideWithPath restores Ev Ion data
func RestoreEVPeptideWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPeptideBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Peptides)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVProteinWithPath restores Ev Protein data
func RestoreEVProteinWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvProteinBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Proteins)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVModsWithPath restores Ev Mods data
func RestoreEVModsWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Mods)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVModificationsWithPath restores Ev Mods data
func RestoreEVModificationsWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsEvBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.Modifications)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}

// RestoreEVCombinedWithPath restores Ev Mods data
func RestoreEVCombinedWithPath(evi *Evidence, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvCombinedBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi.CombinedProtein)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

	return
}
