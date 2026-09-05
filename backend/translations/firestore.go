package translations

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type FirestoreRepository struct {
	client  *firestore.Client
	spaceID string
}

func NewFirestoreRepository(client *firestore.Client, spaceID string) (*FirestoreRepository, error) {
	if client == nil {
		return nil, errors.New("Firestore client is required")
	}
	spaceID = strings.TrimSpace(spaceID)
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	return &FirestoreRepository{client: client, spaceID: spaceID}, nil
}

func (repository *FirestoreRepository) questionRef(questionID string) *firestore.DocumentRef {
	return repository.client.Collection("spaces").Doc(repository.spaceID).Collection("ext").Doc("issuenumber").Collection("questions").Doc(questionID)
}

func (repository *FirestoreRepository) GetQuestion(ctx context.Context, questionID string) (QuestionSource, error) {
	snapshot, err := repository.questionRef(questionID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return QuestionSource{}, ErrQuestionNotFound
		}
		return QuestionSource{}, err
	}
	return decodeQuestion(snapshot)
}

func (repository *FirestoreRepository) GetTranslation(ctx context.Context, questionID, language string) (Translation, error) {
	snapshot, err := repository.questionRef(questionID).Collection("translations").Doc(language).Get(ctx)
	if err != nil {
		return Translation{}, err
	}
	var value Translation
	if err := snapshot.DataTo(&value); err != nil {
		return Translation{}, err
	}
	return value, nil
}

func (repository *FirestoreRepository) PutTranslationIfFresh(ctx context.Context, question QuestionSource, value Translation) (stored Translation, err error) {
	qref := repository.questionRef(question.ID)
	tref := qref.Collection("translations").Doc(value.Language)
	err = repository.client.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {
		snapshot, err := transaction.Get(qref)
		if err != nil {
			return err
		}
		live, err := decodeQuestion(snapshot)
		if err != nil {
			return err
		}
		sourceLanguage := live.SourceLanguage
		if sourceLanguage == "" {
			sourceLanguage = "en"
		}
		if live.ContentRevision != question.ContentRevision || SourceHash(live, sourceLanguage) != value.SourceHash {
			return ErrStaleSource
		}
		if existing, err := transaction.Get(tref); err == nil {
			if err := existing.DataTo(&stored); err != nil {
				return err
			}
			if stored.CorrectedByUID != "" && stored.SourceRevision == value.SourceRevision && stored.SourceHash == value.SourceHash {
				return nil
			}
			if !stored.CreatedAt.IsZero() {
				value.CreatedAt = stored.CreatedAt
			}
		} else if status.Code(err) != codes.NotFound {
			return err
		}
		if value.CreatedAt.IsZero() {
			value.CreatedAt = time.Now().UTC()
		}
		stored = value
		return transaction.Set(tref, stored)
	})
	return stored, err
}

func decodeQuestion(snapshot *firestore.DocumentSnapshot) (QuestionSource, error) {
	var stored struct {
		Title           string `firestore:"title"`
		Description     string `firestore:"description"`
		SourceLanguage  string `firestore:"sourceLanguage"`
		ContentRevision int64  `firestore:"contentRevision"`
		Publication     string `firestore:"publication"`
		Status          string `firestore:"status"`
		CreatorUID      string `firestore:"creatorUID"`
	}
	if err := snapshot.DataTo(&stored); err != nil {
		return QuestionSource{}, fmt.Errorf("decode question %s: %w", snapshot.Ref.ID, err)
	}
	publication := stored.Publication
	if publication == "" {
		publication = stored.Status
	}
	return QuestionSource{ID: snapshot.Ref.ID, Title: stored.Title, Description: stored.Description, SourceLanguage: stored.SourceLanguage, ContentRevision: stored.ContentRevision, Publication: publication, CreatorUID: stored.CreatorUID}, nil
}
