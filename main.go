package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func main() {
	// ----------- ServiceAccount Credential のパスを反映 start --------------
	opt := option.WithCredentialsFile("./credential.json")
	// ----------- ServiceAccount Credential のパスを反映 end --------------

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// ----------- 対象UID設定 start --------------
	uid := "6ndlpCNF5bfQpu5PwHO6vpZW1Up2"
	// ----------- 対象UID設定 end --------------

	// ----------- 各種機能 start（利用するところをコメントアウト） --------------
	// カスタムクレームへのセット
	customClaimsSet(ctx, app, uid)

	// カスタムクレームの確認
	// client, err := app.Auth(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// customClaimsRead(ctx, client, uid)

	// リフレッシュトークンのrevoke
	// revokeRefreshTokens(ctx, app, uid)
	// ----------- 各種機能 start（利用するところをコメントアウト） --------------
}

func customClaimsSet(ctx context.Context, app *firebase.App, uid string) {

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// ----------- ロール設定 変更箇所 start --------------
	claims := map[string]interface{}{
		"role_A": true,
		"role_B": true,
		"role_C": false,
	}
	// ----------- ロール設定 変更箇所 end --------------

	err = client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		log.Fatalf("error setting custom claims %v\n", err)
	}
}

func customClaimsRead(ctx context.Context, client *auth.Client, uid string) {
	// [START read_custom_user_claims_golang]
	// Lookup the user associated with the specified uid.
	user, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatal(err)
	}
	// The claims can be accessed on the user record.
	log.Println(user.CustomClaims)
	// [END read_custom_user_claims_golang]
}

func revokeRefreshTokens(ctx context.Context, app *firebase.App, uid string) {
	// [START revoke_tokens_golang]
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	if err := client.RevokeRefreshTokens(ctx, uid); err != nil {
		log.Fatalf("error revoking tokens for user: %v, %v\n", uid, err)
	}
	// accessing the user's TokenValidAfter
	u, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	timestamp := u.TokensValidAfterMillis / 1000
	log.Printf("the refresh tokens were revoked at: %d (UTC seconds) ", timestamp)
	// [END revoke_tokens_golang]
}
