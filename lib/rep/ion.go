package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"philosopher/lib/msg"

	"philosopher/lib/bio"
	"philosopher/lib/cla"
	"philosopher/lib/id"
	"philosopher/lib/mod"
	"philosopher/lib/sys"
	"philosopher/lib/uti"
)

// AssembleIonReport reports consist on ion reporting
func (evi *Evidence) AssembleIonReport(ion id.PepIDList, decoyTag string) {

	var list IonEvidenceList
	var psmPtMap = make(map[string][]string)
	var psmIonMap = make(map[string][]string)
	var bestProb = make(map[string]float64)

	var ionMods = make(map[string][]mod.Modification)

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range evi.PSM {

		psmIonMap[i.IonForm] = append(psmIonMap[i.IonForm], i.Spectrum)
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.Protein)

		if i.Probability > bestProb[i.IonForm] {
			bestProb[i.IonForm] = i.Probability
		}

		for j := range i.MappedProteins {
			psmPtMap[i.IonForm] = append(psmPtMap[i.IonForm], j)
		}

		for _, j := range i.Modifications.Index {
			ionMods[i.IonForm] = append(ionMods[i.IonForm], j)
		}

	}

	for _, i := range ion {
		var pr IonEvidence

		pr.IonForm = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		pr.Spectra = make(map[string]int)
		pr.MappedGenes = make(map[string]int)
		pr.MappedProteins = make(map[string]int)
		pr.Modifications.Index = make(map[string]mod.Modification)

		v, ok := psmIonMap[pr.IonForm]
		if ok {
			for _, j := range v {
				pr.Spectra[j]++
			}
		}

		pr.Sequence = i.Peptide
		pr.ModifiedSequence = i.ModifiedPeptide
		pr.MZ = uti.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
		pr.ChargeState = i.AssumedCharge
		pr.PeptideMass = i.CalcNeutralPepMass
		pr.PrecursorNeutralMass = i.PrecursorNeutralMass
		pr.Expectation = i.Expectation
		pr.NumberOfEnzymaticTermini = i.NumberOfEnzymaticTermini
		pr.Protein = i.Protein
		pr.MappedProteins[i.Protein] = 0
		pr.Modifications = i.Modifications
		pr.Probability = bestProb[pr.IonForm]

		// get the mapped proteins
		for _, j := range psmPtMap[pr.IonForm] {
			pr.MappedProteins[j] = 0
		}

		mods, ok := ionMods[pr.IonForm]
		if ok {
			for _, j := range mods {
				_, okMod := pr.Modifications.Index[j.Index]
				if !okMod {
					pr.Modifications.Index[j.Index] = j
				}
			}
		}

		// is this bservation a decoy ?
		if cla.IsDecoyPSM(i, decoyTag) {
			pr.IsDecoy = true
		}

		list = append(list, pr)
	}

	sort.Sort(list)
	evi.Ions = list

	return
}

// MetaIonReport reports consist on ion reporting
func (evi Evidence) MetaIonReport(brand string, channels int, hasDecoys bool) {

	var header string
	output := fmt.Sprintf("%s%sion.tsv", sys.MetaDir(), string(filepath.Separator))

	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(errors.New("peptide ion output file"), "fatal")
	}
	defer file.Close()

	// building the printing set tat may or not contain decoys
	var printSet IonEvidenceList
	for _, i := range evi.Ions {
		// This inclusion is necessary to avoid unexistent observations from being included after using the filter --mods options
		if i.Probability > 0 {
			if hasDecoys == false {
				if i.IsDecoy == false {
					printSet = append(printSet, i)
				}
			} else {
				printSet = append(printSet, i)
			}
		}
	}

	header = "Peptide Sequence\tModified Sequence\tPeptide Length\tM/Z\tCharge\tObserved Mass\tProbability\tExpectation\tSpectral Count\tIntensity\tAssigned Modifications\tObserved Modifications\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Genes\tMapped Proteins"

	if brand == "tmt" {
		switch channels {
		case 6:
			header += "\tChannel 126\tChannel 127N\tChannel 128C\tChannel 129N\tChannel 130C\tChannel 131"
		case 10:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N"
		case 11:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C"
		case 16:
			header += "\tChannel 126\tChannel 127N\tChannel 127C\tChannel 128N\tChannel 128C\tChannel 129N\tChannel 129C\tChannel 130N\tChannel 130C\tChannel 131N\tChannel 131C\tChannel 132N\tChannel 132C\tChannel 133N\tChannel 133C\tChannel 134N"
		default:
			header += ""
		}
	} else if brand == "itraq" {
		switch channels {
		case 4:
			header += "\tChannel 114\tChannel 115\tChannel 116\tChannel 117"
		case 8:
			header += "\tChannel 113\tChannel 114\tChannel 115\tChannel 116\tChannel 117\tChannel 118\tChannel 119\tChannel 121"
		default:
			header += ""
		}
	}

	header += "\n"

	// verify if the structure has labels, if so, replace the original channel names by them.
	if len(printSet[0].Labels.Channel1.CustomName) > 3 {
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel1.Name, printSet[0].Labels.Channel1.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel2.Name, printSet[0].Labels.Channel2.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel3.Name, printSet[0].Labels.Channel3.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel4.Name, printSet[0].Labels.Channel4.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel5.Name, printSet[0].Labels.Channel5.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel6.Name, printSet[0].Labels.Channel6.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel7.Name, printSet[0].Labels.Channel7.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel8.Name, printSet[0].Labels.Channel8.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel9.Name, printSet[0].Labels.Channel9.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel10.Name, printSet[0].Labels.Channel10.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel11.Name, printSet[0].Labels.Channel11.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel12.Name, printSet[0].Labels.Channel12.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel13.Name, printSet[0].Labels.Channel13.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel14.Name, printSet[0].Labels.Channel14.CustomName, -1)
		header = strings.Replace(header, "Channel "+printSet[0].Labels.Channel15.Name, printSet[0].Labels.Channel15.CustomName, -1)
	}

	_, e = io.WriteString(file, header)
	if e != nil {
		msg.WriteToFile(errors.New("Cannot print Ion to file"), "fatal")
	}

	for _, i := range printSet {

		assL, obs := getModsList(i.Modifications.Index)

		var mappedProteins []string
		for j := range i.MappedProteins {
			if j != i.Protein {
				mappedProteins = append(mappedProteins, j)
			}
		}

		var mappedGenes []string
		for j := range i.MappedGenes {
			if j != i.GeneName && len(j) > 0 {
				mappedGenes = append(mappedGenes, j)
			}
		}

		sort.Strings(mappedGenes)
		sort.Strings(mappedProteins)
		sort.Strings(assL)
		sort.Strings(obs)

		line := fmt.Sprintf("%s\t%s\t%d\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%.4f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			i.Sequence,
			i.ModifiedSequence,
			len(i.Sequence),
			i.MZ,
			i.ChargeState,
			i.PeptideMass,
			i.Probability,
			i.Expectation,
			len(i.Spectra),
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.Protein,
			i.ProteinID,
			i.EntryName,
			i.GeneName,
			i.ProteinDescription,
			strings.Join(mappedGenes, ","),
			strings.Join(mappedProteins, ","),
		)

		switch channels {
		case 4:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
			)
		case 6:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
			)
		case 8:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
				i.Labels.Channel7.Intensity,
				i.Labels.Channel8.Intensity,
			)
		case 10:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
				i.Labels.Channel7.Intensity,
				i.Labels.Channel8.Intensity,
				i.Labels.Channel9.Intensity,
				i.Labels.Channel10.Intensity,
			)
		case 11:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
				i.Labels.Channel7.Intensity,
				i.Labels.Channel8.Intensity,
				i.Labels.Channel9.Intensity,
				i.Labels.Channel10.Intensity,
				i.Labels.Channel11.Intensity,
			)
		case 16:
			line = fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f",
				line,
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
				i.Labels.Channel7.Intensity,
				i.Labels.Channel8.Intensity,
				i.Labels.Channel9.Intensity,
				i.Labels.Channel10.Intensity,
				i.Labels.Channel11.Intensity,
				i.Labels.Channel12.Intensity,
				i.Labels.Channel13.Intensity,
				i.Labels.Channel14.Intensity,
				i.Labels.Channel15.Intensity,
				i.Labels.Channel16.Intensity,
			)
		default:
			header += ""
		}

		line += "\n"

		_, e = io.WriteString(file, line)
		if e != nil {
			msg.WriteToFile(errors.New("Cannot print Ions to file"), "fatal")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}
