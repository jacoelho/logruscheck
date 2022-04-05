package logs

import (
	log "github.com/sirupsen/logrus"
)

func LogWithFields() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}

func LogWithFieldsNotInSnakeCase() {
	log.WithFields(log.Fields{
		"marineAnimal": "walrus", // want "log key 'marineAnimal' should be snake_case."
	}).Info("A walrus appears")
}

func LogWithFieldsMapStringInterface() {
	log.WithFields(map[string]interface{}{
		"bar": "bar",
	}).Debug("something")
}

func LogWithFieldsMapStringInterfaceNotInSnakeCase() {
	log.WithFields(map[string]interface{}{
		"barBar": "bar", // want "log key 'barBar' should be snake_case."
	}).Debug("something")
}

func LogWithFieldsMapStringAny() {
	log.WithFields(map[string]any{
		"foo": "foo",
	}).Debug("something")
}

func LogWithFieldsMapStringAnyNotInSnakeCase() {
	log.WithFields(map[string]any{
		"fooFoo": "foo", // want "log key 'fooFoo' should be snake_case."
	}).Debug("something")
}

func LogWithField() {
	log.WithField("foo", "bar").Warn("foo")
}

func LogWithFieldNotInSnakeCase() {
	log.
		WithField("fooFoo", "bar"). // want "log key 'fooFoo' should be snake_case."
		Warn("foo")
}

func LogTracef() {
	log.Tracef("foo") // want "call to 'Tracef' should be replaced with WithField or WithFields."
}

func LogDebugf() {
	log.Debugf("foo") // want "call to 'Debugf' should be replaced with WithField or WithFields."
}

func LogPrintf() {
	log.Printf("foo") // want "call to 'Printf' should be replaced with WithField or WithFields."
}

func LogWarnf() {
	log.Warnf("foo") // want "call to 'Warnf' should be replaced with WithField or WithFields."
}

func LogWarningf() {
	log.Warningf("foo") // want "call to 'Warningf' should be replaced with WithField or WithFields."
}

func LogErrorf() {
	log.Errorf("foo") // want "call to 'Errorf' should be replaced with WithField or WithFields."
}

func LogPanicf() {
	log.Panicf("foo") // want "call to 'Panicf' should be replaced with WithField or WithFields."
}

func LogFatalf() {
	log.Fatalf("foo") // want "call to 'Fatalf' should be replaced with WithField or WithFields."
}
