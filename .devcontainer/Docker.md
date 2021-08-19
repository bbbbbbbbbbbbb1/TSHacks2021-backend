# Dockerを用いて開発環境をセットアップする
1. Visual Studio Codeに拡張機能「Remote - Containers」をインストール
2. 右下の`><`を押して「Reopen in Container」
3. 自動的にコンテナイメージのビルドやコンテナ作成が実行される
4. 実行に必要な`dlv-dap`をインストールする
  * コマンドパレット(Ctrl(MacはCommand) + Shift + P)から「Go: Install/Update Tools」を実行してインストール
  * 開発するなら全てインストールした方が良い

2021/08/19時点での最新はfeature/herokuであり，localhost:8080で繋がる．

あとは普通にデバッグも可能．
コンソールはコンテナ内のものであることには注意する必要がある(gitのconfigなど)．

# TODO
[ ] mysqlのコンテナ作成