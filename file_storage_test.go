package starrpg

import (
	"testing"
)

func TestGet(t *testing.T) {
	/*
	 * パス毎に GET 結果、 PUT 可否 (フォーマットチェック)、 DELETE 可否が異なる。
	 * /<model>
	 *   POST
	 * /<model>/<id>
	 *   PUT, DELETE (親子を考慮して消さないといけない)
	 * /<model>/<id>/<property>
	 *   PUT
	 * /<parent-model>/<parent-id>/<model>
	 *   POST
	 *
	 * /games
	 *   [{id: 1, name: "Piyo"}, {...}, ...]
	 * /games/12342342
	 *   {id: 12342342, name: "Foo Foo", maps: [12345, 12346]}
	 * /games/12342342/name
	 *   "foo"
	 * /games/12342342/maps
	 *   [12345, 12346]? 名前は? ラベルとして必要?
	 * /maps
	 *   ?
	 * /maps/12345
	 *   {}
	 * /planes/1234473294879
	 *   {id: 1234473294879, value: "..."}
	 */
}
