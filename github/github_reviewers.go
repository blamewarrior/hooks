package github

type Reviewers interface {
	AssignReviewers(repositoryFullName string, reviewers []Collaborator) (err error)
	ReviewComments(commentsURl string) (string, error)
}
