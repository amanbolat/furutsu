import Vue from 'vue'
import {
    Dialog,
    Popover,
    Loading,
    MessageBox,
    Message,
    Notification,
    InputNumber
} from 'element-ui'

Vue.use(Loading.directive)
Vue.prototype.$loading = Loading.service
Vue.prototype.$msgbox = MessageBox
Vue.prototype.$alert = MessageBox.alert
Vue.prototype.$confirm = MessageBox.confirm
Vue.prototype.$prompt = MessageBox.prompt
Vue.prototype.$notify = Notification
Vue.prototype.$message = Message
Vue.use(Dialog)
Vue.use(InputNumber)
Vue.use(Popover)
