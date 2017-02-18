package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/spec"
)

func ValidateCompositionById(compositionId string) (Data, error) {
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	composition := bigchain.GetTxData(tx)
	if err = spec.ValidComposition(composition); err != nil {
		return nil, err
	}
	return ValidateComposition(composition, pub)
}

func ValidateComposition(composition Data, pub crypto.PublicKey) (Data, error) {
	composerId := spec.GetCompositionComposer(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	composerPub := bigchain.GetTxPublicKey(tx)
	composer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(composer); err != nil {
		return nil, err
	}
	publisherId := spec.GetCompositionPublisher(composition)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	publisherPub := bigchain.GetTxPublicKey(tx)
	publisher := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(publisher); err != nil {
		return nil, err
	}
	if pub.Equals(composerPub) {
		//..
	} else if pub.Equals(publisherPub) {
		//..
	} else {
		return nil, ErrInvalidKey
	}
	rights := spec.GetCompositionRights(composition)
	for _, right := range rights {
		rightHolderId := spec.GetRightHolder(right)
		tx, err = bigchain.GetTx(rightHolderId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rightHolder := bigchain.GetTxData(tx)
		if err = spec.ValidAgent(rightHolder); err != nil {
			return nil, err
		}
	}
	return composition, nil
}

func ValidateRecordingById(recordingId string) (Data, error) {
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	recording := bigchain.GetTxData(tx)
	if err = spec.ValidRecording(recording); err != nil {
		return nil, err
	}
	return ValidateRecording(recording, pub)
}

func ValidateRecording(recording Data, pub crypto.PublicKey) (Data, error) {
	compositionId := spec.GetRecordingComposition(recording)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return nil, err
	}
	rights := spec.GetCompositionRights(composition)
	rightHolderIds := make(map[string]struct{})
	for _, right := range rights {
		rightHolderIds[spec.GetRightHolder(right)] = struct{}{}
	}
	labelId := spec.GetRecordingLabel(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	labelPub := bigchain.GetTxPublicKey(tx)
	label := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(label); err != nil {
		return nil, err
	}
	performerId := spec.GetRecordingPerformer(recording)
	tx, err = bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	performerPub := bigchain.GetTxPublicKey(tx)
	performer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(performer); err != nil {
		return nil, err
	}
	producerId := spec.GetRecordingProducer(recording)
	tx, err = bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	producer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(producer); err != nil {
		return nil, ErrorAppend(ErrInvalidModel, spec.AGENT)
	}
	if pub.Equals(labelPub) {
		if _, ok := rightHolderIds[labelId]; !ok {
			return nil, ErrorAppend(ErrCriteriaNotMet, "label is not composition right holder")
		}
	} else if pub.Equals(performerPub) {
		if _, ok := rightHolderIds[performerId]; !ok {
			return nil, ErrorAppend(ErrCriteriaNotMet, "performer is not composition right holder")
		}
	} else {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recording must be signed by label or performer")
	}
	rights = spec.GetRecordingRights(recording)
	for _, right := range rights {
		rightHolderId := spec.GetRightHolder(right)
		tx, err = bigchain.GetTx(rightHolderId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rightHolder := bigchain.GetTxData(tx)
		if err = spec.ValidAgent(rightHolder); err != nil {
			return nil, ErrorAppend(ErrInvalidModel, spec.AGENT)
		}
	}
	return recording, nil
}

func ValidatePublishingLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	license := bigchain.GetTxData(tx)
	if err = spec.ValidPublishingLicense(license); err != nil {
		return nil, err
	}
	return ValidatePublishingLicense(license, pub)
}

func ValidatePublishingLicense(license Data, pub crypto.PublicKey) (Data, error) {
	compositionId := spec.GetLicenseComposition(license)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return nil, err
	}
	rights := spec.GetCompositionRights(composition)
	rightHolderIds := make(map[string]struct{})
	for _, right := range rights {
		rightHolderIds[spec.GetRightHolder(right)] = struct{}{}
	}
	licenserId := spec.GetLicenseLicenser(license)
	if _, ok := rightHolderIds[licenserId]; !ok {
		return nil, ErrorAppend(ErrCriteriaNotMet, "licenser is not a composition right holder")
	}
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrInvalidKey
	}
	licenseeId := spec.GetLicenseLicensee(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return nil, err
	}
	return license, nil
}

func ValidateRecordingLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	license := bigchain.GetTxData(tx)
	if err = spec.ValidRecordingLicense(license); err != nil {
		return nil, err
	}
	return ValidateRecordingLicense(license, pub)
}

func ValidateRecordingLicense(license Data, pub crypto.PublicKey) (Data, error) {
	recordingId := spec.GetLicenseRecording(license)
	recording, err := ValidateRecordingById(recordingId)
	if err != nil {
		return nil, err
	}
	rights := spec.GetRecordingRights(recording)
	rightHolderIds := make(map[string]struct{})
	for _, right := range rights {
		rightHolderIds[spec.GetRightHolder(right)] = struct{}{}
	}
	licenserId := spec.GetLicenseLicenser(license)
	if _, ok := rightHolderIds[licenserId]; !ok {
		return nil, ErrorAppend(ErrCriteriaNotMet, "licenser is not a recording right holder")
	}
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrInvalidKey
	}
	licenseeId := spec.GetLicenseLicensee(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return nil, err
	}
	return license, nil
}