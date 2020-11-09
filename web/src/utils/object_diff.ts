import transform from 'lodash/transform'
import isEqual from 'lodash/isEqual'
import isObject from 'lodash/isObject'

/**
 * Deep diff between two object, using lodash
 */
function difference(object: any, base: any, omitKeys?: string[]): any {
    function changes(obj: any, bas: any) {
        return transform(obj, (result: any, value, key: string) => {
            if (omitKeys) {
                if (!omitKeys.includes(key)) {
                    if (!isEqual(value, bas[key])) {
                        result[key] = (isObject(value) && isObject(bas[key])) ? changes(value, bas[key]) : value;
                    }
                }
            } else {
                if (!isEqual(value, bas[key])) {
                    result[key] = (isObject(value) && isObject(bas[key])) ? changes(value, bas[key]) : value;
                }
            }
        });
    }
    return changes(object, base);
}

export default difference