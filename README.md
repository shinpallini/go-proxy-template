# go-proxy-template

Golangを使ったproxyサーバーのサンプルです。

以下の3つのコンポーネントで構成されています

1. apiserver
    - clientがアクセスしたいapi server
2. proxy
    - apiserverのresponseを中継したいserver
3. client
    - apiserverの情報をproxy経由で取得するclient