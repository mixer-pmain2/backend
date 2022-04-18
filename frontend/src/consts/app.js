import {accessModule} from "../configs/access";

export const app = {
  DISPANSER: 0,
  STAC: 1
}

export const dispanserSubModules = {
  visit: {
    id: 0,
    title: "Посещения",
    img: require("assets/images/statoscop.svg"),
    unit: accessModule.dispanser.visit.unit,
    access: accessModule.dispanser.visit.access,
  },
  history: {
    id: 1,
    title: "История",
    unit: accessModule.dispanser.history.unit,
    access: accessModule.dispanser.history.access,
  },
  uchet: {
    id: 2,
    title: "Учет",
    unit: accessModule.dispanser.uchet.unit,
    access: accessModule.dispanser.uchet.access,
  },
  invalid: {
    id: 3,
    title: "Инвалидность",
    unit: accessModule.dispanser.invalid.unit,
    access: accessModule.dispanser.invalid.access,
  },
  custody: {
    id: 4,
    title: "Опекунство",
    unit: accessModule.dispanser.custody.unit,
    access: accessModule.dispanser.custody.access,
  },
  groupWork: {
    id: 5,
    title: "Групповая работа",
    unit: accessModule.dispanser.groupWork.unit,
    access: accessModule.dispanser.groupWork.access,
  },
  vaccinations: {
    id: 6,
    title: "Прививки",
    unit: accessModule.dispanser.vaccinations.unit,
    access: accessModule.dispanser.vaccinations.access,
  },
  infection: {
    id: 7,
    title: "Инфекционные заболевания",
    unit: accessModule.dispanser.infection.unit,
    access: accessModule.dispanser.infection.access,
  },
  UKL: {
    id: 8,
    title: "УКЛ",
    unit: accessModule.dispanser.UKL.unit,
    access: accessModule.dispanser.UKL.access,
  },
  prof: {
    id: 9,
    title: "Проф.",
    unit: accessModule.dispanser.Prof.unit,
    access: accessModule.dispanser.Prof.access,
  }

}