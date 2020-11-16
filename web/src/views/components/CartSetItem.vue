<template lang="pug">
  v-list-item(v-if="item && item.product")
    v-list-item-avatar
      v-img(:src="`/img/${item.product.name}.jpg`")
    v-list-item-content
      v-btn(icon @click="deleteItem(item)").delete-item-btn
        v-icon(color="red") mdi-trash-can-outline
      v-list-item-subtitle.text-capitalize.black--text {{ item.product.name }}
      v-list-item-subtitle
        span.black--text Price: {{ sumFromCents(item.product.price).toFixed(2)  }}$
          v-edit-dialog.mt-2(
            :return-value.sync="item"
            large
            @open="openAmountDialog"
            @save="saveAmountDialog"
            @close="closeAmountDialog"
            @cancel="closeAmountDialog")
            .black--text.font-weight-bold.mr-2.amount-text.d-flex Amount: {{ item.amount }}
              v-icon.ml-1(small) mdi-square-edit-outline
            template(v-slot:input)
              NumberInput(:disabled="loading" :value="item.amount" :min="1" :step="1" @change="handleAmountChange(item, $event)" inline controls size="small")
</template>

<script lang="ts">
import {Component, Emit, Mixins, Prop, Watch, PropSync} from 'vue-property-decorator'
import AppMixin from '../../mixins/AppMixin'
import _isNumber from 'lodash/isNumber'

@Component({
  name: 'CartSetItem',
})
export default class CartItem extends Mixins(AppMixin) {
  @Prop(Object) public cartItem!: any
  @PropSync('items', {type: Object}) public cartItems!: any[]
  @PropSync('isLoading', {type: Boolean}) public loading!: any[]
  public item: any = {}
  public amountSaved = false

  @Watch('cartItem')
  onCartItemChange(val: any) {
    if (!val) {
      return
    }
    this.item = JSON.parse(JSON.stringify(val))
  }

  public created() {
    if (!this.cartItem) {
      return
    }
    this.item = JSON.parse(JSON.stringify(this.cartItem))
  }

  @Emit()
  private deleteItem(item: any) {
    return item
  }

  @Emit()
  public changeItemAmount(item: any, amount: number) {
    return {item, amount}
  }

  public handleAmountChange(item: any, amount: number) {
    if (!_isNumber(amount)) {
      return
    }
    if (amount < 1) {
      return
    }
    item.amount = amount
  }

  private openAmountDialog() {
    this.amountSaved = false
  }

  private closeAmountDialog() {
    if (this.amountSaved) {
      return
    }
    this.item.amount = this.cartItem.amount
  }

  private saveAmountDialog() {
    if (this.item.amount === this.cartItem.amount) {
      return
    }
    const diff = this.cartItems[this.item.product.id].amount + (this.item.amount - this.cartItem.amount)
    this.changeItemAmount(this.item, diff)
    this.amountSaved = true
  }

  private calcItemTotal(item: any): number {
    return (item.product.price / 100 * item.amount)
  }

  private calcItemDiscountedTotal(item: any): number {
    let total = item.product.price / 100 * item.amount
    total = total - (total * item.discount_percent / 100)
    return total
  }
}
</script>
