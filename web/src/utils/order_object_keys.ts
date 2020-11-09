import sort from 'lodash/sortBy'

export default function (cart: any): any {
    const copy = JSON.parse(JSON.stringify(cart))
    const nonDiscSet = sort(copy.non_discount_set, function (o: any) {
        return o.id
    })

    const newDiscSet: any[] = []
    if (copy.discount_sets) {
        copy.discount_sets.forEach((val: any) => {
            const set = sort(val, function (o: any) {
                return o.id
            })
            newDiscSet.push(set)
        })
    }


    copy.non_discount_set = nonDiscSet
    copy.discount_sets = newDiscSet

    return copy
}