package tokenize

import v1beta2 "cloud.google.com/go/language/apiv1beta2/languagepb"

type (
	PartOfSpeechTag         int32
	PartOfSpeechAspect      int32
	PartOfSpeechCase        int32
	PartOfSpeechForm        int32
	PartOfSpeechGender      int32
	PartOfSpeechMood        int32
	PartOfSpeechNumber      int32
	PartOfSpeechPerson      int32
	PartOfSpeechProper      int32
	PartOfSpeechReciprocity int32
	PartOfSpeechTense       int32
	PartOfSpeechVoice       int32
	PartOfSpeech            struct {
		Tag         PartOfSpeechTag
		Aspect      PartOfSpeechAspect
		Case        PartOfSpeechCase
		Form        PartOfSpeechForm
		Gender      PartOfSpeechGender
		Mood        PartOfSpeechMood
		Number      PartOfSpeechNumber
		Person      PartOfSpeechPerson
		Proper      PartOfSpeechProper
		Reciprocity PartOfSpeechReciprocity
		Tense       PartOfSpeechTense
		Voice       PartOfSpeechVoice
	}
)

const (
	// Tags
	PartOfSpeechTagUnknown = PartOfSpeechTag(v1beta2.PartOfSpeech_UNKNOWN)
	PartOfSpeechTagAdj     = PartOfSpeechTag(v1beta2.PartOfSpeech_ADJ)
	PartOfSpeechTagAdp     = PartOfSpeechTag(v1beta2.PartOfSpeech_ADP)
	PartOfSpeechTagAdv     = PartOfSpeechTag(v1beta2.PartOfSpeech_ADV)
	PartOfSpeechTagConj    = PartOfSpeechTag(v1beta2.PartOfSpeech_CONJ)
	PartOfSpeechTagDet     = PartOfSpeechTag(v1beta2.PartOfSpeech_DET)
	PartOfSpeechTagNoun    = PartOfSpeechTag(v1beta2.PartOfSpeech_NOUN)
	PartOfSpeechTagNum     = PartOfSpeechTag(v1beta2.PartOfSpeech_NUM)
	PartOfSpeechTagPron    = PartOfSpeechTag(v1beta2.PartOfSpeech_PRON)
	PartOfSpeechTagPrt     = PartOfSpeechTag(v1beta2.PartOfSpeech_PRT)
	PartOfSpeechTagPunct   = PartOfSpeechTag(v1beta2.PartOfSpeech_PUNCT)
	PartOfSpeechTagVerb    = PartOfSpeechTag(v1beta2.PartOfSpeech_VERB)
	PartOfSpeechTagX       = PartOfSpeechTag(v1beta2.PartOfSpeech_X)
	PartOfSpeechTagAffix   = PartOfSpeechTag(v1beta2.PartOfSpeech_AFFIX)
	// Aspect
	PartOfSpechAspectUnknown      = PartOfSpeechAspect(v1beta2.PartOfSpeech_ASPECT_UNKNOWN)
	PartOfSpechAspectPerfective   = PartOfSpeechAspect(v1beta2.PartOfSpeech_PERFECTIVE)
	PartOfSpechAspectImperfective = PartOfSpeechAspect(v1beta2.PartOfSpeech_IMPERFECTIVE)
	PartOfSpechAspectProgressive  = PartOfSpeechAspect(v1beta2.PartOfSpeech_PROGRESSIVE)
	// Case
	PartOfSpeechCaseUnknown       = PartOfSpeechCase(v1beta2.PartOfSpeech_CASE_UNKNOWN)
	PartOfSpeechCaseAccusative    = PartOfSpeechCase(v1beta2.PartOfSpeech_ACCUSATIVE)
	PartOfSpeechCaseAdverbial     = PartOfSpeechCase(v1beta2.PartOfSpeech_ADVERBIAL)
	PartOfSpeechCaseComplemantive = PartOfSpeechCase(v1beta2.PartOfSpeech_COMPLEMENTIVE)
	PartOfSpeechCaseDative        = PartOfSpeechCase(v1beta2.PartOfSpeech_DATIVE)
	PartOfSpeechCaseGenitive      = PartOfSpeechCase(v1beta2.PartOfSpeech_GENITIVE)
	PartOfSpeechCaseInstrumental  = PartOfSpeechCase(v1beta2.PartOfSpeech_INSTRUMENTAL)
	PartOfSpeechCaseLocative      = PartOfSpeechCase(v1beta2.PartOfSpeech_LOCATIVE)
	PartOfSpeechCaseNominative    = PartOfSpeechCase(v1beta2.PartOfSpeech_NOMINATIVE)
	PartOfSpeechCaseOblique       = PartOfSpeechCase(v1beta2.PartOfSpeech_OBLIQUE)
	PartOfSpeechCasePartitive     = PartOfSpeechCase(v1beta2.PartOfSpeech_PARTITIVE)
	PartOfSpeechCasePrepositional = PartOfSpeechCase(v1beta2.PartOfSpeech_PREPOSITIONAL)
	PartOfSpeechCaseReflexive     = PartOfSpeechCase(v1beta2.PartOfSpeech_REFLEXIVE_CASE)
	PartOfSpeechCaseRelative      = PartOfSpeechCase(v1beta2.PartOfSpeech_RELATIVE_CASE)
	PartOfSpeechCaseVocative      = PartOfSpeechCase(v1beta2.PartOfSpeech_VOCATIVE)
	// Form
	PartOfSpeechFormUnknown        = PartOfSpeechForm(v1beta2.PartOfSpeech_FORM_UNKNOWN)
	PartOfSpeechFormAdnomial       = PartOfSpeechForm(v1beta2.PartOfSpeech_ADNOMIAL)
	PartOfSpeechFormAuxiliary      = PartOfSpeechForm(v1beta2.PartOfSpeech_AUXILIARY)
	PartOfSpeechFormComplementizer = PartOfSpeechForm(v1beta2.PartOfSpeech_COMPLEMENTIZER)
	PartOfSpeechFormFinalEnding    = PartOfSpeechForm(v1beta2.PartOfSpeech_FINAL_ENDING)
	PartOfSpeechFormGerund         = PartOfSpeechForm(v1beta2.PartOfSpeech_GERUND)
	PartOfSpeechFormRealis         = PartOfSpeechForm(v1beta2.PartOfSpeech_REALIS)
	PartOfSpeechFormIrrealis       = PartOfSpeechForm(v1beta2.PartOfSpeech_IRREALIS)
	PartOfSpeechFormShort          = PartOfSpeechForm(v1beta2.PartOfSpeech_SHORT)
	PartOfSpeechFormLong           = PartOfSpeechForm(v1beta2.PartOfSpeech_LONG)
	PartOfSpeechFormOrder          = PartOfSpeechForm(v1beta2.PartOfSpeech_ORDER)
	PartOfSpeechFormSpecific       = PartOfSpeechForm(v1beta2.PartOfSpeech_SPECIFIC)
	// Gender
	PartOfSpeechGenderUnknown   = PartOfSpeechGender(v1beta2.PartOfSpeech_GENDER_UNKNOWN)
	PartOfSpeechGenderFeminine  = PartOfSpeechGender(v1beta2.PartOfSpeech_FEMININE)
	PartOfSpeechGenderMasculine = PartOfSpeechGender(v1beta2.PartOfSpeech_MASCULINE)
	PartOfSpeechGenderNeuter    = PartOfSpeechGender(v1beta2.PartOfSpeech_NEUTER)
	// Mood
	PartOfSpeechMoodUnknown       = PartOfSpeechMood(v1beta2.PartOfSpeech_MOOD_UNKNOWN)
	PartOfSpeechMoodConditional   = PartOfSpeechMood(v1beta2.PartOfSpeech_CONDITIONAL_MOOD)
	PartOfSpeechMoodImperative    = PartOfSpeechMood(v1beta2.PartOfSpeech_IMPERATIVE)
	PartOfSpeechMoodIndicative    = PartOfSpeechMood(v1beta2.PartOfSpeech_INDICATIVE)
	PartOfSpeechMoodInterrogative = PartOfSpeechMood(v1beta2.PartOfSpeech_INTERROGATIVE)
	PartOfSpeechMoodJussive       = PartOfSpeechMood(v1beta2.PartOfSpeech_JUSSIVE)
	PartOfSpeechMoodSubjunctive   = PartOfSpeechMood(v1beta2.PartOfSpeech_SUBJUNCTIVE)
	// Number
	PartOfSpeechNumberUnknown  = PartOfSpeechNumber(v1beta2.PartOfSpeech_NUMBER_UNKNOWN)
	PartOfSpeechNumberSingular = PartOfSpeechNumber(v1beta2.PartOfSpeech_SINGULAR)
	PartOfSpeechNumberPlural   = PartOfSpeechNumber(v1beta2.PartOfSpeech_PLURAL)
	PartOfSpeechNumberDual     = PartOfSpeechNumber(v1beta2.PartOfSpeech_DUAL)
	// Person
	PartOfSpeechPersonUnknown   = PartOfSpeechPerson(v1beta2.PartOfSpeech_PERSON_UNKNOWN)
	PartOfSpeechPersonFirst     = PartOfSpeechPerson(v1beta2.PartOfSpeech_FIRST)
	PartOfSpeechPersonSecond    = PartOfSpeechPerson(v1beta2.PartOfSpeech_SECOND)
	PartOfSpeechPersonThird     = PartOfSpeechPerson(v1beta2.PartOfSpeech_THIRD)
	PartOfSpeechPersonReflexive = PartOfSpeechPerson(v1beta2.PartOfSpeech_REFLEXIVE_PERSON)
	// Proper
	PartOfSpeechProperUnknown = PartOfSpeechProper(v1beta2.PartOfSpeech_PROPER_UNKNOWN)
	PartOfSpeechIsProper      = PartOfSpeechProper(v1beta2.PartOfSpeech_PROPER)
	PartOfSpeechIsNotProper   = PartOfSpeechProper(v1beta2.PartOfSpeech_NOT_PROPER)
	// Reciprocity
	PartOfSpeechReciprocityUnknown       = PartOfSpeechReciprocity(v1beta2.PartOfSpeech_RECIPROCITY_UNKNOWN)
	PartOfSpeechReciprocityReciprocal    = PartOfSpeechReciprocity(v1beta2.PartOfSpeech_RECIPROCAL)
	PartOfSpeechReciprocityNonReciprocal = PartOfSpeechReciprocity(v1beta2.PartOfSpeech_NON_RECIPROCAL)
	// Tense
	PartOfSpeechTenseUnknown     = PartOfSpeechTense(v1beta2.PartOfSpeech_TENSE_UNKNOWN)
	PartOfSpeechTenseConditional = PartOfSpeechTense(v1beta2.PartOfSpeech_CONDITIONAL_TENSE)
	PartOfSpeechTenseFuture      = PartOfSpeechTense(v1beta2.PartOfSpeech_FUTURE)
	PartOfSpeechTensePast        = PartOfSpeechTense(v1beta2.PartOfSpeech_PAST)
	PartOfSpeechTensePresent     = PartOfSpeechTense(v1beta2.PartOfSpeech_PRESENT)
	PartOfSpeechTenseImperfect   = PartOfSpeechTense(v1beta2.PartOfSpeech_IMPERFECT)
	PartOfSpeechTensePluperfect  = PartOfSpeechTense(v1beta2.PartOfSpeech_PLUPERFECT)
	// Voice
	PartOfSpeechVoiceUnknown   = PartOfSpeechVoice(v1beta2.PartOfSpeech_VOICE_UNKNOWN)
	PartOfSpeechVoiceActive    = PartOfSpeechVoice(v1beta2.PartOfSpeech_ACTIVE)
	PartOfSpeechVoiceCausative = PartOfSpeechVoice(v1beta2.PartOfSpeech_CAUSATIVE)
	PartOfSpeechVoicePassive   = PartOfSpeechVoice(v1beta2.PartOfSpeech_PASSIVE)
)
