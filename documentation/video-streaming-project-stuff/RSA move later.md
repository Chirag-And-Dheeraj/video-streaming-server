# RSA

- choose 2 [large] prime numbers p & q
- calculate N [N = pq]
- calculate T [T = (p-1)(q-1)]   [Euler Totient]
- choose 2 numbers e & d | (ed%T) = 1 && e < T && e must be coprime with N and T
- publish N & e [public key]
- keep d secret [private key]

## Example
p = 2
q = 5
N = 2 * 5 = 10
T = (2 - 1) * (5 - 1) = 4


Building the program

- areCoPrime(x,y)  ..................  x2 (e with N, e with T)
- 