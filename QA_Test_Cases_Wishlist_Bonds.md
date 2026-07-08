# QA Test Cases — Wishlist & Bonds Module

## 1. Wishlist Name Limiting

- [ ] 1. What if the user enters emojis in the wishlist name?
- [ ] 2. What if the wishlist name starts with an underscore?
- [ ] 3. What if the wishlist name contains underscores in the middle/end?
- [ ] 4. What if the wishlist name contains special characters (e.g. @, #, $, %, &)?
- [ ] 5. What if the wishlist name contains numbers?
- [ ] 6. What if the wishlist name contains spaces between words?
- [ ] 7. What if the wishlist name starts with a space?
- [ ] 8. What if the wishlist name consists only of spaces?
- [ ] 9. What if the wishlist name field is left blank/empty?
- [ ] 10. What if the wishlist name exceeds the maximum allowed character limit?
- [ ] 11. What if the user enters a wishlist name that is only 1 character long?

## 2. Max Limit of Wishlists (>5 shows warning)

- [ ] 12. Can the user create more than 5 wishlists?
- [ ] 13. What if the user creates exactly 5 wishlists, does the warning appear prematurely?
- [ ] 14. What if the user tries to create a 6th wishlist, does the warning show correctly?
- [ ] 15. Can the user still create a wishlist after the warning is shown, or is creation blocked?
- [ ] 16. What happens if the user deletes a wishlist after hitting the limit, can they create a new one again?

## 3. Duplicate Bonds Not Allowed in the Same Group

- [ ] 17. What if the user tries to add the same bond twice to the same group?
- [ ] 18. Can the user add the same bond to two different groups?
- [ ] 19. What message/feedback is shown when a duplicate bond addition is attempted?

## 4. Max 10 Bonds Allowed in a Group

- [ ] 20. Can the user add more than 10 bonds in a wishlist/group?
- [ ] 21. What happens when the user tries to add the 11th bond to a group?
- [ ] 22. Is a warning/error message shown when the bond limit is reached?

## 5. Adding a Sold-Out Bond

- [ ] 23. What happens if a user tries to add a bond that is sold out?
- [ ] 24. Is the 'Add to Wishlist' option disabled or hidden for sold-out bonds?

## 6. Multiple Clicks / Expired Bond Handling

- [ ] 25. What happens if a user clicks multiple times rapidly on 'Add to Wishlist' for a bond?
- [ ] 26. What happens if a user clicks on an expired bond to add it to the wishlist?
- [ ] 27. Does rapid multiple clicking create duplicate entries in the database?

## 7. Duplicate Group Names

- [ ] 28. What happens if a user creates a group with a name that already exists?
- [ ] 29. What happens if a user renames a group to a name that already exists?
- [ ] 30. Is the duplicate check case-sensitive (e.g. 'Group1' vs 'group1')?

## 8. Adding a Near-Expiring Bond

- [ ] 31. What happens if a user tries to add a near-expiring bond to the wishlist?
- [ ] 32. Is any warning/indicator shown for near-expiring bonds before adding?

## 9. Removing All Bonds from a Group

- [ ] 33. What happens when a user removes every bond from a group?
- [ ] 34. Does the group get deleted automatically once it becomes empty?
- [ ] 35. Or does the group remain in the system as an empty group?

## 10. Add to Wishlist Database Record Behavior

- [ ] 36. If a user clicks 'Add to Wishlist' once, does it create a single record in the DB?
- [ ] 37. Or does a single click create multiple records in the DB?

## 11. Search Functionality

- [ ] 38. What happens if the user searches for something invalid or nonsensical (e.g. 'trash' input)?
- [ ] 39. What if the user searches with special characters or emojis?
- [ ] 40. What if the user searches with an empty/blank query?

## 12. Multi-Device / Concurrency Handling

- [ ] 41. What happens if a user tries to add a bond at the same time on two different devices logged into the same account?
- [ ] 42. Does this create duplicate entries or a conflict error?

## 13. Group Deletion Behavior

- [ ] 43. If a user deletes a group, does it also delete the bonds/items inside that group?
- [ ] 44. In the DB, is the group hard-deleted (permanently removed)?
- [ ] 45. Or is it soft-deleted (marked inactive but retained)?
- [ ] 46. What if the user tries to undo a group deletion, is recovery possible?
