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


### Benchmark Result

* env:
    * Node: MacBook Pro (Retina, 13-inch, Early 2015)
    * CPU: 2.9 GHz Intel Core i5
    * MEM: 16 GB 1867 MHz DDR3

* program:
    * go-himeno: ebc904ad7462778e05949adf229e075ae80b2544

| MFLOPS                     | SSMALL   | SMALL    | MIDDLE   | LARGE    | ELARGE   |
| -------------------------- | -------- | -------- | -------- | -------- | -------- |
| go-himeno (Thread: 4)      | 3919.153 | 4092.738 | 4039.370 | 4065.742 | 1951.904 |
| go-himeno (Thread: 1)      | 2036.409 | 2258.040 | 2246.051 | 2150.445 | 1053.351 |
| original-C-himeno (CPU: 1) | 3928.642 | 3667.226 | 3764.491 | 3594.881 | 2905.840 |

