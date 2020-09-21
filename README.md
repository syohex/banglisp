# banglisp

Simple LISP implementation of [対話によるCommon Lisp入門](https://www.morikita.co.jp/books/book/2093) in Golang

## Build

```bash
% cd cmd/banglisp
% go build
```

## Run REPL

```bash
% ./banglisp
> (+ 1 2)
3
> (sin (/ pi 2))
1E+00
> (cons "foo" (cons "bar" nil))
(foo bar)
```

## Run file

```lisp
(defun fizzbuzz1 (m n)
  (if (= m n)
      nil
    (if (and (= (mod m 3) 0) (= (mod m 5) 0))
	(cons "fizzbuzz" (fizzbuzz1 (+ m 1) n))
      (if (= (mod m 5) 0)
	  (cons "buzz" (fizzbuzz1 (+ m 1) n))
	(if (= (mod m 3) 0)
	    (cons "fizz" (fizzbuzz1 (+ m 1) n))
	  (cons m (fizzbuzz1 (+ m 1) n)))))))

(defun fizzbuzz (n)
  (fizzbuzz1 1 n))

(print (fizzbuzz 16))
```

```bash
% ./banglisp fizzbuzz.lisp
(1 2 "fizz" 4 "buzz" "fizz" 7 8 "fizz" "buzz" 11 "fizz" 13 14 "fizzbuzz")
```