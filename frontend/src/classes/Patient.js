export default class Patient {
    constructor(objPatient) {
        this.patient = objPatient
    }

    getLastUchet() {
        if (this.patient?.uchet?.length > 0)
            return this.patient.uchet[0]
        return {}
    }

}