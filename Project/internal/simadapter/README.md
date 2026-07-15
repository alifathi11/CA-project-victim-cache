# Akita adapter layer

این پوشه محل اتصال مدل مستقل پروژه به Akita است.

TODOهای مرحله اتصال:

- ساخت Engine آکیتا
- تعریف Component/Spec/State برای traffic generator، L1، Victim Cache، L2 و Memory
- تعریف پورت‌های ورودی و خروجی
- ساخت connectionها
- تبدیل `model.Request` به پیام Akita
- زمان‌بندی latencyها با زمان شبیه‌سازی
- ثبت trace و metric

در مرحله اول هیچ import از Akita وجود ندارد تا اسکلت به API یک نسخه خاص قفل نشود.
