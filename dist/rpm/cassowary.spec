%global bindir /usr/local/sbin

Name: cassowary
Version: 0.14.0
Release: 1
Summary: Modern cross-platform HTTP load-testing tool written in Go
License: MIT
Group: Development/Tools
URL: https://github.com/rogerwelin/cassowary
Source0: https://github.com/rogerwelin/cassowary/releases/download/v%{version}/cassowary_Linux_x86_64.tar.gz

%description
Modern cross-platform HTTP load-testing tool written in Go

%prep
%setup -cqn cassowary_%{version}_Linux_x86_64

%install
%{__mkdir_p} %{buildroot}/%{bindir}
%{__cp} -a cassowary %{buildroot}/%{bindir}/.

%clean
%{__rm} -rf %{buildroot}

%files
%doc README.md LICENSE
%defattr(-,root,root,-)
%{bindir}/cassowary

%changelog
* Fri Jan 24 2020 <roger.welin@icloud.com>
- Cassowary spec file
* Mon Jun 13 2022 <ericchou19831101@msn.com>
- Uplift release version
