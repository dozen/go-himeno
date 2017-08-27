## HOW TO USE

### compile

```
git clone github.com/dozen/go-himeno
cd go-himeno
make SIZE=MIDDLE
```

SIZE: [SSMALL, SMALL, MIDDLE, LARGE, ELARGE]


### run

```
./go-himeno 
```

change worker goroutine

```
./go-himeno 8 8
```

最初の引数はJacobiのルーチンのワーカー数、2個目は配列のコピーをするルーチンのワーカー数。

デフォルトではどちらもCPU数になる。


## Benchmark Result

* env:
    * Node: MacBook Pro (Retina, 13-inch, Early 2015)
    * CPU: 2.9 GHz Intel Core i5
    * MEM: 16 GB 1867 MHz DDR3

* program:
    * go-himeno: ebc904ad7462778e05949adf229e075ae80b2544
    * go-himeno (ELARGE): 02980e19556fa3a5e42e0d545b33b24d2cf865ed

| # | MFLOPS                     | SSMALL   | SMALL    | MIDDLE   | LARGE    | ELARGE   |
| - | -------------------------- | -------- | -------- | -------- | -------- | -------- |
| 1 | go-himeno (Worker: 4)      | 3919.153 | 4092.738 | 4039.370 | 4065.742 | 1951.904 |
| 2 | go-himeno (Worker: 1)      | 2036.409 | 2258.040 | 2246.051 | 2150.445 | 1053.351 |
| 3 | original-C-himeno (CPU: 1) | 3928.642 | 3667.226 | 3764.491 | 3594.881 | 2905.840 |
|   | #1 slower than C-himeno    | x1.00    | x0.90    | x0.93    | x0.88    | x1.49    |
|   | #2 slower than C-himeno    | x1.93    | x1.62    | x1.68    | x1.67    | x2.76    |

#### 備考
ELARGEのベンチマークのみパッチを当てた_himeno.goを使用している。配列のポインタを宣言するよう変更しているので、多少遅くなる。
手元の環境ではgo-himenoをELARGEで動かすとmacOSのメモリ圧縮機能によってメモリが圧縮されてしまったので、ベンチ結果に影響が出ていると思われる。
