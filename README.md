# C-Style JavaScript

A "small" project I am messing around with because I had the idea of making a C-Style language that compiles to fully type-safe JavaScript.

### Example:

Input:

```
int add(int a, int b) {
	return a + b;
}
```

Output:

```js
function add(a, b) {
	if (typeof a !== "number" || !Number.isInteger(a)) {
		throw new TypeError("add: param a must be an integer");
	}

	if (typeof b !== "number" || !Number.isInteger(b)) {
		throw new TypeError("add: param b must be an integer");
	}

	return a + b;
}
```
