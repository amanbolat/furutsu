import {ActionContext, ActionTree, GetterTree, MutationTree} from 'vuex'
import Auth from '../api/auth'
import {User} from '../api/auth'
import {decodeJwtToken} from '../utils/jwt'

export class State {
    public user: User = {} as User
    public authToken: string = ''
}

const getters = {
    user(state: State): User {
        return state.user
    },
    authToken(state: State): string {
        return state.authToken
    },
    isAuthenticated(state: State): boolean {
        return state.authToken !== ''
    }
} as GetterTree<State, any>

const mutations = {
    SET_USER(state: State, usr: User) {
        state.user = usr
    },
    SET_AUTH_TOKEN(state: State, token: string) {
        state.authToken = token
    },
    DEL_AUTH_TOKEN(state: State) {
        state.authToken = ''
    }
} as MutationTree<State>

const actions = {
    Login(store: ActionContext<State, any>, credentials: any) {
        return new Promise((resolve, reject) => {
            Auth.login(credentials.username, credentials.password).then((response: any) => {
                console.log("AUTH", response)
                const usr = decodeJwtToken(response) as User
                console.log(usr)
                store.commit('SET_USER', usr)
                store.commit('SET_AUTH_TOKEN', response)

                resolve()
            }).catch((error: any) => {
                reject(error)
            })
        })
    },

    Logout(store: ActionContext<State, any>) {
        return new Promise((resolve) => {
            store.commit('DEL_AUTH_TOKEN', '')
            store.commit('SET_USER', {})
            resolve()
        })
    }
} as ActionTree<State, any>

const user = {
    state: new State(),
    getters,
    mutations,
    actions
}

export default user
