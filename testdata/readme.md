
# Testdata generation

```bash
sqlite3 -csv map.sqlite "select pos, hex(data) from blocks"
```