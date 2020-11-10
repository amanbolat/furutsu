<template lang="pug">
  v-list-item(v-if="item && item.product")
    v-list-item-avatar
      v-img(:src="`/img/${item.product.name}.jpg`")
    v-list-item-content
      v-list-item-subtitle.text-capitalize.black--text {{ item.product.name }}
      v-list-item-subtitle Price: {{ sumFromCents(item.product.price).toFixed(2) }}$
        a-input-number.ml-5(:value="item.amount" :min="0" :step="1" size="small" @change="changeItemAmount(item, $event)")
</template>

<script lang="ts">
import {Component, Emit, Mixins, Prop, Watch, PropSync} from 'vue-property-decorator'
import AppMixin from '../../mixins/AppMixin'

@Component({
  name: 'CartSetItem',
})
export default class CartItem extends Mixins(AppMixin) {
  @Prop(Object) public cartItem!: any
  @PropSync('items', {type: Object}) public cartItems!: any[]
  private item = {}

  @Watch('cartItem')
  onCartItemChange(val: any) {
    if (!val) {
      return
    }
    this.item = JSON.parse(JSON.stringify(val))
  }

  private created() {
    if (!this.cartItem) {
      return
    }
    this.item = JSON.parse(JSON.stringify(this.cartItem))
  }

  @Emit()
  public changeItemAmount(item: any, amount: number) {
    const diff = this.cartItems[item.product.id].amount + (amount - this.cartItem.amount)
    item.amount = amount

    return {item, diff}
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
