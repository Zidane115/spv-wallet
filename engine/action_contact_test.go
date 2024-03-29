package engine

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	xPubGeneric            = "62910a1ecbc7123213231563ab3f8aa70568ed934d1e0383cb1bbbfb1bc8f2afe5"
	xPubForNotFoundContact = "7fa312762ef940d9f744906913422d750e76b980e5824cc7995a2d803af765ee2c"
	paymailGeneric         = "test@test.test"
	fullName               = "John Doe"
	pubKey                 = "pubKey"
)

type testCase struct {
	testID               int
	name                 string
	data                 testCaseData
	expectedErrorMessage error
}

type testCaseData struct {
	xPub          string
	paymail       string
	contactStatus string
	deleted       bool
}

func initContactTestCase(t *testing.T) (context.Context, ClientInterface, func()) {
	ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())

	xPub := newXpub(testXPub, append(client.DefaultModelOptions(), New())...)
	err := xPub.Save(ctx)
	require.NoError(t, err)

	return ctx, client, deferMe
}

func TestAcceptContactHappyPath(t *testing.T) {
	ctx, client, deferMe := initContactTestCase(t)
	defer deferMe()

	t.Run("accept contact, should return nil", func(t *testing.T) {
		// given
		contact := newContact(
			fullName,
			paymailGeneric,
			pubKey,
			xPubGeneric,
			ContactAwaitAccept,
		)
		contact.enrich(ModelContact, append(client.DefaultModelOptions(), New())...)
		err := contact.Save(ctx)
		require.NoError(t, err)

		// when
		err = client.AcceptContact(ctx, xPubGeneric, paymailGeneric)

		// then
		require.NoError(t, err)

		contact1, err := getContact(ctx, paymailGeneric, xPubGeneric, client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.Equal(t, ContactNotConfirmed, contact1.Status)
	})
}

func TestAcceptContactErrorPath(t *testing.T) {

	testCases := []testCase{
		{
			testID: 1,
			name:   "non existance contact, should return \"contact not found\" error",
			data: testCaseData{
				xPub:          xPubForNotFoundContact,
				paymail:       paymailGeneric,
				contactStatus: ContactAwaitAccept.String(),
			},
			expectedErrorMessage: ErrContactNotFound,
		},
		{
			testID: 2,
			name:   "contact has status \"confirmed\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactAwaitAccept.String(),
			},
			expectedErrorMessage: ErrContactIncorrectStatus,
		},
		{
			testID: 3,
			name:   "contact has status \"not confirmed\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactNotConfirmed.String(),
			},
			expectedErrorMessage: ErrContactIncorrectStatus,
		},
		{
			testID: 4,
			name:   "contact has status \"rejected\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactRejected.String(),
			},
			expectedErrorMessage: ErrContactIncorrectStatus,
		},
		{
			testID: 5,
			name:   "contact has status \"rejected\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactRejected.String(),
				deleted:       true,
			},
			expectedErrorMessage: ErrContactNotFound,
		},
	}

	for _, tc := range testCases {
		ctx, client, deferMe := initContactTestCase(t)
		defer deferMe()
		t.Run(tc.name, func(t *testing.T) {
			// given
			contact := newContact(
				fullName,
				paymailGeneric,
				pubKey,
				xPubGeneric,
				ContactNotConfirmed,
			)
			contact.enrich(ModelContact, append(client.DefaultModelOptions(), New())...)
			if tc.data.deleted {
				contact.DeletedAt.Valid = true
				contact.DeletedAt.Time = time.Now()
			}
			err := contact.Save(ctx)
			require.NoError(t, err)

			// when
			err = client.AcceptContact(ctx, tc.data.xPub, tc.data.paymail)

			// then
			require.Error(t, err)
			require.EqualError(t, err, tc.expectedErrorMessage.Error())
		})
	}
}

func TestRejectContactHappyPath(t *testing.T) {
	ctx, client, deferMe := initContactTestCase(t)
	defer deferMe()

	t.Run("reject contact", func(t *testing.T) {
		// given
		contact := newContact(
			fullName,
			paymailGeneric,
			pubKey,
			xPubGeneric,
			ContactAwaitAccept,
		)
		contact.enrich(ModelContact, append(client.DefaultModelOptions(), New())...)
		err := contact.Save(ctx)
		require.NoError(t, err)

		// when
		err = client.RejectContact(ctx, contact.OwnerXpubID, contact.Paymail)

		// then
		require.NoError(t, err)

		contact1, err := getContact(ctx, contact.Paymail, contact.OwnerXpubID, client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.Empty(t, contact1)
	})
}

func TestRejectContactErrorPath(t *testing.T) {

	testCases := []testCase{
		{
			testID: 1,
			name:   "non existance contact, should return \"contact not found\" error",
			data: testCaseData{
				xPub:          xPubForNotFoundContact,
				paymail:       paymailGeneric,
				contactStatus: ContactAwaitAccept.String(),
			},
			expectedErrorMessage: ErrContactNotFound,
		},
		{
			testID: 2,
			name:   "contact has status \"confirmed\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactConfirmed.String(),
			},
			expectedErrorMessage: ErrContactIncorrectStatus,
		},
		{
			testID: 3,
			name:   "contact has status \"not confirmed\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactNotConfirmed.String(),
			},
			expectedErrorMessage: ErrContactIncorrectStatus,
		},
		{
			testID: 4,
			name:   "contact has status \"rejected\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactRejected.String(),
			},
			expectedErrorMessage: ErrContactIncorrectStatus,
		},
		{
			testID: 5,
			name:   "contact has status \"rejected\", should return \"contact does not have status awaiting\" error",
			data: testCaseData{
				xPub:          xPubGeneric,
				paymail:       paymailGeneric,
				contactStatus: ContactRejected.String(),
				deleted:       true,
			},
			expectedErrorMessage: ErrContactNotFound,
		}}

	for _, tc := range testCases {
		ctx, client, deferMe := initContactTestCase(t)
		defer deferMe()
		t.Run(tc.name, func(t *testing.T) {
			// given
			contact := newContact(
				fullName,
				paymailGeneric,
				pubKey,
				xPubGeneric,
				ContactNotConfirmed,
			)
			contact.enrich(ModelContact, append(client.DefaultModelOptions(), New())...)
			if tc.data.deleted {
				contact.DeletedAt.Valid = true
				contact.DeletedAt.Time = time.Now()
			}
			err := contact.Save(ctx)
			require.NoError(t, err)

			// when
			err = client.RejectContact(ctx, tc.data.xPub, tc.data.paymail)

			// then
			require.Error(t, err)
			require.EqualError(t, err, tc.expectedErrorMessage.Error())
		})
	}
}
