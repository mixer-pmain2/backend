import {formatDateToInput} from "./string";

export default (params) => {

  const getParam = (params, param) => {
    return params.filter(v => v.param === param)?.[0]
  }

  const getMinVisitDate = () => {
    const date = new Date()
    const max_day = getParam(params, "max_day")?.paramI

    if (max_day <= date.getDate()) {
      date.setDate(1)
    } else date.setMonth(date.getMonth() - 1)

    return formatDateToInput(date)
  }

  const getMaxVisitDate = () => {
    const date = new Date()
    return formatDateToInput(date)
  }

  const _params = {
    visit: {
      minDate: getMinVisitDate(),
      maxDate: getMaxVisitDate()
    }
  }


  return _params
}

export const sleep = (milliseconds) => {
  const date = Date.now();
  let currentDate = null;
  do {
    currentDate = Date.now();
  } while (currentDate - date < milliseconds);
}