create trigger updateUserSocialRating
    after insert
    on feedback
    for row
execute procedure updateUserSocialRating();

create trigger validateFeedbackPoints
    after insert
    on feedback
    for statement
execute procedure validateFeedbackPoints();

create trigger validateOrdersPrice
    after insert
    on executor_order
    for statement
execute procedure validateOrdersPrice();

create trigger validateOrdersPrice
    after insert
    on customer_order
    for statement
execute procedure validateOrdersPrice();

create trigger validateCarLicensePlate
    after insert
    on car
    for statement
execute procedure validateCarLicensePlate();
