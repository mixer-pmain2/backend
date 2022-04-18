export const getAge = (bday, toDate) => {
  console.log(bday, toDate)
  let old = -1
  for (let i = bday; i < toDate; i.setFullYear(i.getFullYear() + 1)) {
    old += 1
  }
  return old
}

export const stringToDate = (s) => {
  const split = s.split("-")

  return new Date(`${split[1]}/${split[2]}/${split[0]}`)
}