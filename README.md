# houtu
A scaffold for basic development based on vue and golang. 

# admin 
- admin user account name must be 'admin'. 
- admin role must be 'admin', so change auth.json carefully.
- if admin password lost, then add etc/config.yaml app section's 
adminpassword attribute for your admin password, when app starts, it will change 
  automatically your admin password, and at last remove this attribute, or else
  it will change to this password again when app starts.
- about cache, you need to choose cache type  
- rbac, auth.json specifies vue router display and apis explain these functions could 
be reached. Just write basic_policy.csv and basic_model.conf to support RBAC.
  
# others
* [vue-admin-template](https://github.com/PanJiaChen/vue-admin-template) 
