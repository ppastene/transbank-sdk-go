# "Primera clase" en Go

Cuando en Go se habla de que algo es de **primera clase** (first-class), nos referimos a
que ese algo se trata **como cualquier otro valor**: se le aplican las mismas reglas que a
un entero, un string o un booleano.

La clave del concepto es la **simetría**: lo que puedes hacer con un entero
(guardarlo, pasarlo, devolverlo) también lo puedes hacer con una función.

## Las 3 operaciones de un valor de primera clase

Un valor es de primera clase cuando puedes hacer con él todo esto:

1. **Almacenarlo** — en una variable, un campo de struct, un elemento de slice o map.
2. **Pasarlo como argumento** a una función.
3. **Devolverlo como retorno** de una función.

El lenguaje no distingue "esto es una función, así que solo se declara y se llama".
Todo se trata con las mismas reglas.

## Ejemplo: funciones que reciben y devuelven funciones

```go
func conAumento(inc int) func(int) int {
    // recibo inc... y devuelvo una función
    return func(x int) int { return x + inc }
}

doble := conAumento(2)   // doble es un valor normal
fmt.Println(doble(5))    // 7
fmt.Println(conAumento(10)(5)) // 15
```

Fíjate que `conAumento(10)` *produce* un valor (una función) que se puede usar al
instante, igual que `5 * 2` produce un entero. Esa es la simetría: **la función es un
valor más del sistema de tipos**.

A las funciones que reciben o devuelven otras funciones se les llama
*higher-order functions* (funciones de orden superior), y los cierres (closures) son
un caso clásico de esto.

## Contraste con "segunda clase"

El concepto opuesto es "segunda clase": entidades que el lenguaje trata con reglas
especiales y limitadas. Por ejemplo, en algunos lenguajes las funciones solo se pueden
declarar y llamar, pero no asignar a una variable ni pasarlas como argumento.
En Go no es así.

## No solo funciones: todos los valores son primera clase

En Go esto aplica a todas las entidades, no solo a las funciones:

- Canales: se pueden guardar en slices, pasarse, devolverse (`[]chan int`).
- Mapas y slices: son valores manipulables.
- Interfaces: se guardan en structs, se pasan, se retornan.

Es lo que hace triviales patrones como la **inyección de dependencias** o los
**middlewares** (`middleware(handler)` recibe y devuelve funciones), que es exactamente
el mecanismo que usa este SDK para permitir inyectar un `http.Client` propio en
`transbank.Options.HTTPClient`.
