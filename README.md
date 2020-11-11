# Furutsu - an online shop for selling fruits

## About
This is the project developed by me as a homework assessment for the technical interview, which include creating the ac

## TODO
Like every project there are no limits in improvement. This is the list of features and improvements that I would add


## Questions and Suggestions

1. The requirement `sets must be added to their own cart item entry` leaves a lot of question and without 
further discussion of that makes it difficult to add some new features or change the business logic. 

Questions:
- What should be done to the amount of products which doesn't fit that discount set? For instance, we have 10 
pears and 5 bananas, then we get only two sets comprised of `8 pears and 4 bananas` in total and `2 pears + 1 banana`
with no discount applied in two more entries.
- What if customer added pears and bananas, but wants to checkout only the pears without deleting bananas
from the cart? 

Suggestion: 
From the UI perspective adding the items left outside of the discount set to different entry is possible, but not 
recommended. Therefore, I would like to suggest create only one cart item for the product and just show, that part
of its amount might get some discounts.
There are some examples of good implementation of similar logic. Chinese e-commerce giants 
such as `Taobao` and `JD.com` have more complex discount system and leave the choice to customers. They 
provide `% discounts`, `set discounts`, `minimum order discounts` and customers may choose them by themselves.

2. Adding the set of items has also affected the functionality of quantity adjustment. If we add or remove items from the
set it changes the quantity of that particular product in the cart, BUT what about removing the item? Should we remove all
the items of this kind or remove only the exact quantity of the product in the set?
Currently, when user presses the DELETE button all the items of this kind will be removed from the cart.

Suggestion: reconsider the idea of the set. See question #1.   

3. Given cart rules about discounts would seem to be a little vague. 

Questions: 
- Discounts intersection. What if two or more discounts are applicable to the same set of items?
- What if the coupon intersects with the discount? 

My thoughts:
As I mentioned before giving the choice to the customer is the best option. However, it might depend on the real world
business requirements. Anyway I would be grateful to have a chance to discuss that question, just for curiosity.

4. Logic of the coupon is not very clear. It seems, that coupons should be generated somehow and everyone who has a code
can use them. Also, it is claimed to be applied to the cart, but not the item cart which seems to me as a bad logic. Indeed, 
the implementation of such a logic simpler, but it might create a few problems in long term. In my opinion, we could attach
coupons to the cart items, thus the discount and coupon discount would be used together or, again, as I said before the 
customer can choose one of the variants, if simple discounts and coupons shouldn't be used together due to business requirements.

