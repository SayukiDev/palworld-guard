# Palworld-Guard
シンプルなパルワールドサーバープログラムのデーモンプログラム

## 使い方
1. Releasesで構築済みのパッケージファイルをダウンロードする
2. パッケージファイルをアンパックする
3. オプションファイルである `config.json5` をチェックする
4. プログラムファイルである `palworld-guard` を実行する

## 利用可能な機能

- デーモン関連
  * 定期的の再起動
  * メモリ利用率がある数値超えたときの再起動
  * クラッシュやキルされたときの自動的再起動
- ディスコードボット関連
  * サーバーの再起動
  * 利用率のチェック
  * ログインしてるプレイヤーリストの取得（プレイヤーネームが英語じゃない場合は失敗する）

## 追加する予定のある機能

- デーモン関連
  * 自動バックアップ
- ディスコードボット関連
  * カスタムコマンドの実行
  * セーブデータの切り替えや追加
  * サーバーオプションの設定

## License
```
    Copyright (C) 2024  Sayuki

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
```