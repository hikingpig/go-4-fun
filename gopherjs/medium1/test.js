m = require("./medium1.js")

h = m.hashit("This is a test")

console.log("Test hashit result: ", h.length)

h = m.hashobj({
    name: "Kent",
    address: "Somewhere in Space and Time",
    city: "Specificity",
    state: "CT"
})
console.log("Test hashobj result: ", h.slice(0, 4))

h = m.hashobj2({
  name: "Kent",
  address: "Somewhere in Space and Time",
  city: "Specificity",
  state: "CT"
})
console.log("Test hashobj2 result: ", h.hash.slice(0, 4))
