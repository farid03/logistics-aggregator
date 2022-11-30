create or replace procedure getCars(id integer) as
$$
begin
    select car.id, car.license_plate, car.position_id from car where owner_id = id;
end;
$$ language plpgsql;

create or replace procedure getCustomerOrders(id integer) as
$$
begin
    select * from customer_order where owner_id = id;
end;
$$ language plpgsql;

create or replace procedure getExecutorOrders(id integer) as
$$
begin
    select * from executor_order where owner_id = id;
end;
$$ language plpgsql;

create or replace function updateUserSocialRating() returns trigger as
$sys_user$
begin
    update sys_user
    set social_rating = (select avg(points) from feedback where receiver_id = old.id)
    where sys_user.id = old.id;
end;
$sys_user$ language plpgsql;

create or replace function validateFeedbackPoints() returns trigger as
$feedback$
begin
    if (new.Points < 0) then
        RAISE EXCEPTION 'Points cannot be negative!';
    end if;
    if (new.Points > 5) then
        RAISE EXCEPTION 'Points cannot be more than 5!';
    end if;
    return new;
end;
$feedback$ language plpgsql;

create or replace function validateOrdersPrice() returns trigger as
$orders$
begin
    if (new.Price < 0) then
        RAISE EXCEPTION 'Points cannot be negative!';
    end if;
    return new;
end;
$orders$ language plpgsql;

create or replace function validateCarLicensePlate() returns trigger as
$orders$
begin
    if (new.license_plate in (select license_plate from car)) then
        RAISE EXCEPTION 'License plate must be unique!';
    end if;
    return new;
end;
$orders$ language plpgsql;

-- если пользользователь не доставил товар, то его нужно заблокировать
-- если пользователь не оплатил доставку, то забанить