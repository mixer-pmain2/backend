//сравниваем права Юзера с правами доступа для "Объекта"
export const isAccessed = (accessList, uUnit, uAccess) => {
    for (let i = 0; i < accessList.length; i++) {
        const unit = accessList[i].unit
        const access = accessList[i].access
        if (uUnit === unit && (uAccess & access) > 0)
            return true
    }
    return false
}

export const accessDispPage = [
    {unit: 1, access: 1},
    {unit: 2, access: 1},
    {unit: 4, access: 1},
    {unit: 8, access: 1},
    {unit: 16, access: 1},
    {unit: 1024, access: 1},
    {unit: 2048, access: 1},
    {unit: 16777216, access: 1},
    {unit: 33554432, access: 1},
]

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