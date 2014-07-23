# Kocha  [![Build Status](https://travis-ci.org/naoina/kocha.png?branch=master)](https://travis-ci.org/naoina/kocha)

A convenient web application framework for [Go](http://golang.org/)

## Features

* Batteries included
* All configurations are in Go's syntax
* Generate an All-In-One binary
* Compatible with `net/http`

## Requirement <a id="Requirement"></a>

* Go 1.3 or later

## Getting started

1. install the framework:

        go get -u github.com/naoina/kocha

    And command-line tool

        go get -u github.com/naoina/kocha/kocha

2. Create a new application:

        kocha new myapp

    Where "myapp" is the application name.

3. Change directory and run the application:

        cd myapp
        kocha run

    or

        cd myapp
        go build -o myapp
        ./myapp

## Documentation

See http://naoina.github.io/kocha/ for more information.

## License

Kocha is licensed under the MIT
