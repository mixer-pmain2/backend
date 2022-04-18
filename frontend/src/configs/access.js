//сравниваем права Юзера с правами доступа для "Объекта"
import {Access, Unit} from "../consts/user";

export const isAccessed = (accessList, uUnit, uAccess) => {
  for (let i = 0; i < accessList.length; i++) {
    const unit = accessList[i].unit
    const access = accessList[i].access
    if (uUnit === unit && (uAccess & access) > 0)
      return true
  }
  return false
}


export const AccessRoleASU = Access.dispanser["Прямой доступ к данным"]
export const AccessRoleDoct = Access.dispanser["Прием пациентов"]
export const AccessRoleRegistrator = Access.dispanser["Работа регистратора"] | Access.dispanser["Только просмотр (справочная система)"]
export const AccessRoleUKL = Access.dispanser["Ввод УКЛ 1 уровня"] |
  Access.dispanser["Ввод УКЛ 2 уровня"] |
  Access.dispanser["Ввод УКЛ 3 уровня"]

export const accessDispPage = [
  {unit: 1, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 2, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 4, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 8, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 16, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 1024, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 2048, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 16777216, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
  {unit: 33554432, access: AccessRoleDoct | AccessRoleRegistrator | AccessRoleASU | AccessRoleUKL},
]
const unitDisp = Unit.АПЛ |
  Unit["Взрослая психиатрия"] |
  Unit.Специалисты |
  Unit["Детская консультация"] |
  Unit["Детский диспансер"] |
  Unit["Подростковая психиатрия"] |
  Unit.ОИЛС |
  Unit.Психотерапия |
  Unit.Суицидология

export const accessModule = {
  dispanser: {
    visit: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    history: {
      access: AccessRoleASU | AccessRoleDoct | AccessRoleRegistrator,
      unit: unitDisp,
    },
    uchet: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    invalid: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    custody: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    groupWork: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    vaccinations: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    infection: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: unitDisp,
    },
    UKL: {
      access: AccessRoleASU | AccessRoleUKL,
      unit: unitDisp,
    },
    Prof: {
      access: AccessRoleASU | AccessRoleDoct,
      unit: Unit["Подростковая психиатрия"] | Unit["Детская консультация"] | Unit["Детский диспансер"],
    }
  }
}

//Страница админки (АСУ)
export const accessAdminPage = [
  {unit: 1, access: 1073741824},
  {unit: 2, access: 1073741824},
  {unit: 4, access: 1073741824},
  {unit: 8, access: 1073741824},
  {unit: 16, access: 1073741824},
  {unit: 32, access: 1073741824},
  {unit: 64, access: 1073741824},
  {unit: 128, access: 1073741824},
  {unit: 256, access: 1073741824},
  {unit: 512, access: 1073741824},
  {unit: 1024, access: 1073741824},
  {unit: 2048, access: 1073741824},
  {unit: 4096, access: 1073741824},
  {unit: 8192, access: 1073741824},
  {unit: 65536, access: 1073741824},
  {unit: 131072, access: 1073741824},
  {unit: 262144, access: 1073741824},
  {unit: 524288, access: 1073741824},
  {unit: 1048576, access: 1073741824},
  {unit: 2097152, access: 1073741824},
  {unit: 4194304, access: 1073741824},
  {unit: 8388608, access: 1073741824},
  {unit: 16777216, access: 1073741824},
  {unit: 33554432, access: 1073741824},
  {unit: 67108864, access: 1073741824},
  {unit: 134217728, access: 1073741824},
  {unit: 268435456, access: 1073741824},
  {unit: 536870912, access: 1073741824},
  {unit: 1073741824, access: 1073741824},
]